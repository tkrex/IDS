package processing

import (
	"sync/atomic"
	"sync"
	"fmt"

	"github.com/tkrex/IDS/common/models"
	"github.com/tkrex/IDS/daemon/persistence"
	"time"
	"encoding/json"
	"github.com/tkrex/IDS/common/utilities"
	"math"
)

type TopicProcessor struct {
	// Flag to indicate that the consumer state
	// 0 == Running
	// 1 == Closed
	state                    int64
	processorStarted         sync.WaitGroup
	processorStopped         sync.WaitGroup
	databaseDelegate         *persistence.DomainInformationStorage
	topicUpdates             []*models.RawTopicMessage
	incomingTopicChannel     chan *models.RawTopicMessage

	forwardingSignalChannel  chan int
	newTopicsCounter         int
}

const (
	TopicForwardThreshold = 10
	SimilarityCheckInterval time.Duration = 7 * 24 * time.Hour
	BulkUpdateThreshold = 10
)

func NewTopicProcessor(incomingTopicChannel chan *models.RawTopicMessage) *TopicProcessor {
	processor := new(TopicProcessor)
	processor.processorStarted.Add(1)
	processor.processorStopped.Add(1)
	processor.incomingTopicChannel = incomingTopicChannel
	processor.forwardingSignalChannel = make(chan int)
	processor.topicUpdates = make([]*models.RawTopicMessage, 0, BulkUpdateThreshold)
	go processor.run()
	processor.processorStarted.Wait()
	fmt.Println("Processor Created")
	return processor
}

func (processor *TopicProcessor) State() int64 {
	return atomic.LoadInt64(&processor.state)
}

func (processor *TopicProcessor) ForwardSignalChannel() chan int {
	return processor.forwardingSignalChannel
}

func (processor *TopicProcessor)  Close() {
	fmt.Println("Closing Processor")
	atomic.StoreInt64(&processor.state, 1)
	processor.processorStopped.Wait()
	fmt.Println("Processor Closed")
}

func (processor *TopicProcessor) run() {
	dbDelegate, err := persistence.NewDomainInformationStorage()
	if err != nil {
		fmt.Println("Stopping Topic Processor: No Conbection to DB")
		return
	}

	processor.databaseDelegate = dbDelegate
	defer processor.databaseDelegate.Close()
	processor.processorStarted.Done()
	for closed := atomic.LoadInt64(&processor.state) == 1; !closed; closed = atomic.LoadInt64(&processor.state) == 1 {
		open := processor.checkIncomingTopic()
		if !open {
			processor.Close()
			break
		}
	}
	processor.processorStopped.Done()
}

func (processor *TopicProcessor) checkIncomingTopic() bool {
	rawTopic, ok := <-processor.incomingTopicChannel
	if rawTopic != nil {
		if !containsDataJSON(rawTopic.Payload) {
			fmt.Println("Drop Topic since it does not contain JSON")
			return ok
		}
		processor.collectIncomingTopic(rawTopic)
	}
	return ok
}
func containsDataJSON(data []byte ) bool {
	var jsonData map[string]*json.RawMessage
	err := json.Unmarshal(data,&jsonData)
	if err != nil {
		return false
	}
	return true
}

func (processor *TopicProcessor) collectIncomingTopic(rawTopic *models.RawTopicMessage) {
	processor.topicUpdates = append(processor.topicUpdates, rawTopic)
	if len(processor.topicUpdates) == BulkUpdateThreshold {
		processor.processIncomingTopics()
		processor.updateBrokerStatistics()
		processor.checkForInformationForwarding()

	}
}

func (processor *TopicProcessor) processIncomingTopics() {
	updates := make([]*models.RawTopicMessage, len(processor.topicUpdates))
	copy(updates, processor.topicUpdates)
	processor.topicUpdates = make([]*models.RawTopicMessage, 0, BulkUpdateThreshold)
	fmt.Println("Bulk Update")
	sortedUpdates := processor.sortTopicUpdatesByName(updates)
	topicNames := make([]string, 0, len(sortedUpdates))
	for name, _ := range sortedUpdates {
		topicNames = append(topicNames, name)
	}
	existingTopics, _ := processor.databaseDelegate.FindTopicsByNames(topicNames)
	fmt.Printf("Number of Existing Topics: %d", len(existingTopics))
	processor.mergeExistingTopicsWithUpdates(existingTopics, sortedUpdates)
}


func (processor *TopicProcessor) sortTopicUpdatesByName(topicUpdates []*models.RawTopicMessage) map[string][]*models.RawTopicMessage {
	sortedTopics := make(map[string][]*models.RawTopicMessage)
	for _, topic := range topicUpdates {
		_, entryExists := sortedTopics[topic.Topic]
		if !entryExists {
			sortedTopics[topic.Topic] = make([]*models.RawTopicMessage, 0, len(topicUpdates))

		}
		sortedTopics[topic.Topic] = append(sortedTopics[topic.Topic], topic)
	}
	return sortedTopics
}

func (processor *TopicProcessor) mergeExistingTopicsWithUpdates(existingTopics map[string]*models.TopicInformation, sortedTopics map[string][]*models.RawTopicMessage) {
	resultingTopicUpdates := make([]*models.TopicInformation, 0, len(sortedTopics))
	processor.newTopicsCounter += (len(sortedTopics) - len(existingTopics))

	var brokerDomain *models.RealWorldDomain
	broker, _ := processor.databaseDelegate.FindBroker()
	brokerDomain = broker.RealWorldDomain

	for name, topicArray := range sortedTopics {
		var resultingTopic *models.TopicInformation
		existingTopic, _ := existingTopics[name]
		resultingTopic = existingTopic
		for _, topic := range topicArray {
			resultingTopic = processor.updateTopicInformation(resultingTopic, topic)
		}
		if brokerDomain != nil {
			resultingTopic.Domain = brokerDomain
		}

		resultingTopicUpdates = append(resultingTopicUpdates, resultingTopic)
	}
	err := processor.databaseDelegate.StoreTopics(resultingTopicUpdates)
	if err != nil {
		fmt.Println("could not update topics")
		return
	}
}

func (processor *TopicProcessor) checkForInformationForwarding() {
	if processor.newTopicsCounter >= TopicForwardThreshold {
		fmt.Println("Trigger Forwarding")
		processor.forwardingSignalChannel <- 1
		processor.newTopicsCounter = 0
	}
}

func (processor *TopicProcessor) updateTopicInformation(existingTopic *models.TopicInformation, newTopic *models.RawTopicMessage) *models.TopicInformation {
	var resultingTopic *models.TopicInformation

	if existingTopic == nil {
		resultingTopic = models.NewTopicInformation(newTopic.Topic, string(newTopic.Payload), newTopic.ArrivalTime)
		processor.calculateUpdateBehavior(resultingTopic, 0)
	} else {
		resultingTopic = existingTopic
		newUpdateInterval := int(newTopic.ArrivalTime.Sub(resultingTopic.LastUpdateTimeStamp).Seconds())
		resultingTopic.LastUpdateTimeStamp = newTopic.ArrivalTime
		processor.calculateUpdateBehavior(resultingTopic, newUpdateInterval)
		processor.calculatePayloadSimilarity(resultingTopic, string(newTopic.Payload))
		resultingTopic.LastPayload = string(newTopic.Payload)

	}
	return resultingTopic
}

func (processor *TopicProcessor) calculatePayloadSimilarity(topic *models.TopicInformation, newJSONPayload string) {
	if topic.UpdateBehavior.NumberOfUpdates % 100 == 0 {
		fmt.Println("DEBUG: Similarity Check for Topic: ", topic.Name)
		processor.calculatePayloadSimilarityCheckInterval(topic)

	}
	if topic.UpdateBehavior.NumberOfUpdates  >= 10 && topic.UpdateBehavior.NumberOfUpdates % topic.SimilarityCheckInterval == 0 {
		fmt.Println("calculating similarity")
		oldJsonKeys :=  getKeysFromJSONString(topic.LastPayload)
		newJsonKeys := getKeysFromJSONString(newJSONPayload)
		if len(oldJsonKeys) == 0 || len(newJsonKeys) == 0 {
			fmt.Println("failed to get keys from JSON file")
			return

		}
		hitCounter := 0
		for  _,key := range oldJsonKeys {
			if Include(newJsonKeys,key) {
				hitCounter++
			}
		}

		similarity := float64(hitCounter) / float64(len(oldJsonKeys)) * 100.0
		topic.PayloadSimilarity = utilities.RoundUp(similarity,2)
	}
}

func (processor *TopicProcessor) calculatePayloadSimilarityCheckInterval(topic *models.TopicInformation) {
	topicSimilarityCheckInterval := int(SimilarityCheckInterval.Seconds()) / int(topic.UpdateBehavior.AverageUpdateIntervalInSeconds)
	fmt.Println(topicSimilarityCheckInterval)
	topic.SimilarityCheckInterval = topicSimilarityCheckInterval
}

func Include(array []string, value string) bool {
	for _,element := range array {
		if element == value {
			return true
		}
	}
	return false
}

func getKeysFromJSONString(jsonString string) []string {
	var objmap map[string]*json.RawMessage
	jsonData := []byte(jsonString)
	if err := json.Unmarshal(jsonData, &objmap); err != nil {
		fmt.Println(err)
		return []string{}
	}
	jsonKeys := make([]string, len(objmap))
	for key, _ := range objmap {
		jsonKeys = append(jsonKeys, key)
	}
	fmt.Println(jsonKeys)
	return jsonKeys
}

func (processor *TopicProcessor) updateBrokerStatistics() {
	broker, err := processor.databaseDelegate.FindBroker()
	if err != nil {
		fmt.Println("TOPIC PROCESSOR: ",err)
		return
	}
	numberOfTopics := processor.databaseDelegate.NumberOfTopics()
	broker.Statistics.NumberOfTopics = numberOfTopics
	secondsSinceLastStatisticUpdate := time.Now().Sub(broker.Statistics.LastStatisticUpdate)
	incomingTopicFrequency := BulkUpdateThreshold / secondsSinceLastStatisticUpdate.Seconds()
	broker.Statistics.ReceivedTopicsPerSeconds = incomingTopicFrequency
}

func (processor *TopicProcessor) calculateUpdateBehavior(topic *models.TopicInformation, newUpdateInterval int) {
	updateBehavior := topic.UpdateBehavior
	if updateBehavior.NumberOfUpdates == 0 {
		updateBehavior.NumberOfUpdates++
		return
	} else if topic.UpdateBehavior.NumberOfUpdates == 1 {
		updateBehavior.AverageUpdateIntervalInSeconds = float64(newUpdateInterval)
		updateBehavior.MaximumUpdateIntervalInSeconds = int(newUpdateInterval)
		updateBehavior.MinimumUpdateIntervalInSeconds = int(newUpdateInterval)
	} else if topic.UpdateBehavior.NumberOfUpdates > 1 {
		updateBehavior.MaximumUpdateIntervalInSeconds = utilities.Max(updateBehavior.MaximumUpdateIntervalInSeconds, newUpdateInterval)
		updateBehavior.MinimumUpdateIntervalInSeconds = utilities.Min(updateBehavior.MinimumUpdateIntervalInSeconds, newUpdateInterval)
		updateBehavior.AverageUpdateIntervalInSeconds = (updateBehavior.AverageUpdateIntervalInSeconds * float64(len(updateBehavior.UpdateIntervalsInSeconds)) + float64(newUpdateInterval)) / float64(len(updateBehavior.UpdateIntervalsInSeconds) + 1)
	}

	updateBehavior.UpdateIntervalsInSeconds = append(updateBehavior.UpdateIntervalsInSeconds, float64(newUpdateInterval))
	updateBehavior.NumberOfUpdates++
	updateBehavior.UpdateIntervalDeviation = processor.calculateMeanAbosluteDeviation(updateBehavior)
	processor.determineReliability(topic)

}

func (processor *TopicProcessor) calculateMeanAbosluteDeviation(updateStats *models.UpdateBehavior) float64 {
	if len(updateStats.UpdateIntervalsInSeconds) < 2 {
		return -1.0
	}
	deviationSum := float64(0)
	for _,interval := range updateStats.UpdateIntervalsInSeconds {
		deviation := math.Abs(interval - updateStats.AverageUpdateIntervalInSeconds)
		deviationSum +=  deviation
	}
	meanAbsoluteDeviation := deviationSum / float64(len(updateStats.UpdateIntervalsInSeconds))

	//Reset Array to current average + deviation to avoid memory leak
	if len(updateStats.UpdateIntervalsInSeconds) == 1000 {
		updateStats.UpdateIntervalsInSeconds = make([]float64,0,1000)
		updateStats.UpdateIntervalsInSeconds[0] = updateStats.AverageUpdateIntervalInSeconds + updateStats.UpdateIntervalDeviation
	}
	return meanAbsoluteDeviation
}


func (processor *TopicProcessor)  determineReliability(topic *models.TopicInformation) {
	intervalDeviation := topic.UpdateBehavior.UpdateIntervalDeviation
	reliability := ""
	if intervalDeviation >= 0 && intervalDeviation <= 30 * time.Minute.Seconds() {
		reliability = "Automatic"
	} else if intervalDeviation <= time.Hour.Seconds() {
		reliability = "Semi-Automatic"
	} else if  intervalDeviation > time.Hour.Seconds() {
		reliability = "Non-Deterministic"
	}
	topic.UpdateBehavior.Reliability = reliability
}