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
	SimilarityCheckInterval = 7 * 24 * time.Hour
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
		open := processor.ProcessIncomingTopics()
		if !open {
			processor.Close()
			break
		}
	}
	processor.processorStopped.Done()
}

func (processor *TopicProcessor) ProcessIncomingTopics() bool {
	rawTopic, ok := <-processor.incomingTopicChannel
	if rawTopic != nil {
		if !containsDataJSON(rawTopic.Payload) {
			fmt.Println("Drop Topic since it does not contain JSON")
			return ok
		}
		processor.processIncomingTopic(rawTopic)

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

func (processor *TopicProcessor) processIncomingTopic(rawTopic *models.RawTopicMessage) {
	processor.topicUpdates = append(processor.topicUpdates, rawTopic)
	if len(processor.topicUpdates) == BulkUpdateThreshold {
		updates := make([]*models.RawTopicMessage, len(processor.topicUpdates))
		copy(updates, processor.topicUpdates)
		processor.topicUpdates = make([]*models.RawTopicMessage, 0, BulkUpdateThreshold)
		fmt.Println("Bulk Update")
		sortedUpdates := processor.sortTopicUpdatesByName(updates)
		fetchRequest := make([]string, 0, len(sortedUpdates))
		for name, _ := range sortedUpdates {
			fetchRequest = append(fetchRequest, name)
		}
		existingTopics, _ := processor.databaseDelegate.FindTopicsByName(fetchRequest)
		fmt.Printf("Number of Existing Topics: %d", len(existingTopics))
		processor.processSortedTopics(existingTopics, sortedUpdates)
	}
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

func (processor *TopicProcessor) processSortedTopics(existingTopics map[string]*models.Topic, sortedTopics map[string][]*models.RawTopicMessage) {
	resultingTopicUpdates := make([]*models.Topic, 0, len(sortedTopics))
	processor.newTopicsCounter += (len(sortedTopics) - len(existingTopics))

	var brokerDomain *models.RealWorldDomain
	broker, _ := processor.databaseDelegate.FindBroker()
	brokerDomain = broker.RealWorldDomains[0]

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

	go processor.triggerDomainInformationUpdate()
}

func (processor *TopicProcessor) triggerDomainInformationUpdate() {
	if processor.newTopicsCounter >= TopicForwardThreshold {
		fmt.Println("Trigger Forwarding")
		processor.forwardingSignalChannel <- 1
		processor.newTopicsCounter = 0
	}
}

func (processor *TopicProcessor) updateTopicInformation(existingTopic *models.Topic, newTopic *models.RawTopicMessage) *models.Topic {
	var resultingTopic *models.Topic

	if existingTopic == nil {
		resultingTopic = models.NewTopic(newTopic.Name, newTopic.Payload, newTopic.ArrivalTime)
		processor.calculateUpdateBehavior(resultingTopic, 0)
	} else {
		resultingTopic = existingTopic
		newUpdateInterval := int(newTopic.ArrivalTime.Sub(resultingTopic.LastUpdateTimeStamp).Seconds())
		processor.calculatePayloadSimilarity(resultingTopic, newTopic.Payload)
		resultingTopic.LastPayload = newTopic.Payload
		resultingTopic.LastUpdateTimeStamp = newTopic.ArrivalTime
		processor.calculateUpdateBehavior(resultingTopic, newUpdateInterval)

	}
	return resultingTopic
}

func (processor *TopicProcessor) calculatePayloadSimilarity(topic *models.Topic, newJSONPayload json.RawMessage) {
	if topic.UpdateBehavior.NumberOfUpdates + 1 >= 10 && topic.UpdateBehavior.NumberOfUpdates % topic.SimilarityCheckInterval == 0 {
		fmt.Println("calculating similarity")
		oldJsonKeys :=  getKeysFromJSONObject(topic.LastPayload)
		newJsonKeys := getKeysFromJSONObject(newJSONPayload)
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

func Include(array []string, value string) bool {
	for _,element := range array {
		if element == value {
			return true
		}
	}
	return false
}

func getKeysFromJSONObject(jsonData []byte) []string {

	var objmap map[string]*json.RawMessage
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
	updateBehavior.UpdateReliability[processor.topicReliabilityStrategy.Name()] = processor.topicReliabilityStrategy.Calculate(updateBehavior)
}


