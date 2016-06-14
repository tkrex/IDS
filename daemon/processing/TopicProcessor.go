package processing

import (
	"sync/atomic"
	"sync"
	"fmt"

	"github.com/tkrex/IDS/common"
	"github.com/tkrex/IDS/common/models"
	"github.com/tkrex/IDS/daemon/persistence"
	"time"
	"encoding/json"
)

type TopicProcessor struct {
	// Flag to indicate that the consumer state
	// 0 == Running
	// 1 == Closed
	state                    int64
	processorStarted         sync.WaitGroup
	processorStopped         sync.WaitGroup
	databaseDelegate         *persistence.DaemonDatabaseWorker
	topicUpdates             []*models.RawTopicMessage
	incomingTopicChannel     chan *models.RawTopicMessage

	forwardingSignalChannel  chan int
	topicReliabilityStrategy models.UpdateReliabilityStrategy

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
	//processor.newTopicsCounter = make(map[string]int)
	processor.forwardingSignalChannel = make(chan int)
	processor.topicReliabilityStrategy = models.MeanAbsoluteDeviation{}
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
	dbDelegate, err := persistence.NewDaemonDatabaseWorker()
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
	existingTopics, _ := processor.databaseDelegate.FindTopicsByName(topicNames)
	fmt.Printf("Number of Existing Topics: %d", len(existingTopics))
	processor.mergeExistingTopicsWithUpdates(existingTopics, sortedUpdates)
}


func (processor *TopicProcessor) sortTopicUpdatesByName(topicUpdates []*models.RawTopicMessage) map[string][]*models.RawTopicMessage {
	sortedTopics := make(map[string][]*models.RawTopicMessage)
	for _, topic := range topicUpdates {
		_, entryExists := sortedTopics[topic.Name]
		if !entryExists {
			sortedTopics[topic.Name] = make([]*models.RawTopicMessage, 0, len(topicUpdates))

		}
		sortedTopics[topic.Name] = append(sortedTopics[topic.Name], topic)
	}
	return sortedTopics
}

func (processor *TopicProcessor) mergeExistingTopicsWithUpdates(existingTopics map[string]*models.Topic, sortedTopics map[string][]*models.RawTopicMessage) {
	resultingTopicUpdates := make([]*models.Topic, 0, len(sortedTopics))
	processor.newTopicsCounter += (len(sortedTopics) - len(existingTopics))

	var brokerDomain *models.RealWorldDomain
	broker, _ := processor.databaseDelegate.FindBroker()
	brokerDomain = broker.RealWorldDomain

	for name, topicArray := range sortedTopics {
		var resultingTopic *models.Topic
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
	_, err := processor.databaseDelegate.StoreTopics(resultingTopicUpdates)
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

func (processor *TopicProcessor) updateTopicInformation(existingTopic *models.Topic, newTopic *models.RawTopicMessage) *models.Topic {
	var resultingTopic *models.Topic

	if existingTopic == nil {
		resultingTopic = models.NewTopic(newTopic.Name, string(newTopic.Payload), newTopic.ArrivalTime)
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

func (processor *TopicProcessor) calculatePayloadSimilarity(topic *models.Topic, newJSONPayload string) {
	if topic.UpdateBehavior.NumberOfUpdates % 100 == 0 {
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
		topic.PayloadSimilarity = common.RoundUp(similarity,2)
	}
}

func (processor *TopicProcessor) calculatePayloadSimilarityCheckInterval(topic *models.Topic) {
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
	numberOfTopics := processor.databaseDelegate.CountTopics()
	broker.Statitics.NumberOfTopics = numberOfTopics
	secondsSinceLastStatisticUpdate := time.Now().Sub(broker.Statitics.LastStatisticUpdate)
	incomingTopicFrequency := BulkUpdateThreshold / secondsSinceLastStatisticUpdate.Seconds()
	broker.Statitics.ReceivedTopicsPerSeconds = incomingTopicFrequency
}

func (processor *TopicProcessor) calculateUpdateBehavior(topic *models.Topic, newUpdateInterval int) {
	updateBehavior := topic.UpdateBehavior
	if updateBehavior.NumberOfUpdates == 0 {
		updateBehavior.NumberOfUpdates++
		return
	} else if topic.UpdateBehavior.NumberOfUpdates == 1 {
		updateBehavior.AverageUpdateIntervalInSeconds = float64(newUpdateInterval)
		updateBehavior.MaximumUpdateIntervalInSeconds = int(newUpdateInterval)
		updateBehavior.MinimumUpdateIntervalInSeconds = int(newUpdateInterval)
	} else if topic.UpdateBehavior.NumberOfUpdates > 1 {
		updateBehavior.MaximumUpdateIntervalInSeconds = common.Max(updateBehavior.MaximumUpdateIntervalInSeconds, newUpdateInterval)
		updateBehavior.MinimumUpdateIntervalInSeconds = common.Min(updateBehavior.MinimumUpdateIntervalInSeconds, newUpdateInterval)
		updateBehavior.AverageUpdateIntervalInSeconds = (updateBehavior.AverageUpdateIntervalInSeconds * float64(len(updateBehavior.UpdateIntervalsInSeconds)) + float64(newUpdateInterval)) / float64(len(updateBehavior.UpdateIntervalsInSeconds) + 1)
	}

	updateBehavior.UpdateIntervalsInSeconds = append(updateBehavior.UpdateIntervalsInSeconds, float64(newUpdateInterval))
	updateBehavior.NumberOfUpdates++
	updateBehavior.UpdateIntervalDeviation = processor.topicReliabilityStrategy.Calculate(updateBehavior)
}