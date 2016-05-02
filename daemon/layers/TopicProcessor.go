package layers

import (
	"sync/atomic"
	"sync"
	"fmt"
	"github.com/tkrex/IDS/common/models"
)

type TopicProcessor struct {

	// Flag to indicate that the consumer state
	// 0 == Running
	// 1 == Closed
	state                int64
	processorStarted     sync.WaitGroup
	processorStopped     sync.WaitGroup
	topicUpdates         []*models.RawTopicMessage
	incomingTopicChannel chan *models.RawTopicMessage
	bulkUpdateThreshold int

}

func NewTopicProcessor(incomingTopicChannel chan *models.RawTopicMessage) *TopicProcessor {
	processor := new(TopicProcessor)
	processor.processorStarted.Add(1)
	processor.processorStopped.Add(1)
	processor.incomingTopicChannel = incomingTopicChannel
	processor.bulkUpdateThreshold = 100
	processor.topicUpdates = make([]*models.RawTopicMessage,0,processor.bulkUpdateThreshold)
	go processor.run()
	processor.processorStarted.Wait()
	fmt.Println("Producer Created")
	return processor
}

func (processor *TopicProcessor) State() int64 {
	return atomic.LoadInt64(&processor.state)
}

func (processor *TopicProcessor)  Close() {
	fmt.Println("Closing Processor")
	atomic.StoreInt64(&processor.state,1)
	processor.processorStopped.Wait()
	fmt.Println("Processor Closed")
}

func  (processor *TopicProcessor) run() {

	processor.processorStarted.Done()
	for closed := atomic.LoadInt64(&processor.state) == 1; !closed; closed = atomic.LoadInt64(&processor.state) == 1 {
		open := processor.ProcessIncomingTopics()
		if !open{
			processor.Close()
			break
		}
	}
	processor.processorStopped.Done()
}

func (processor *TopicProcessor) ProcessIncomingTopics() bool {
	rawTopic, ok :=  <-processor.incomingTopicChannel
	if rawTopic != nil {
		processor.processIncomingTopic(rawTopic)

	}
	return ok
}


func (processor *TopicProcessor) processIncomingTopic(rawTopic *models.RawTopicMessage) {
	processor.topicUpdates = append(processor.topicUpdates,rawTopic)
	if len(processor.topicUpdates) == processor.bulkUpdateThreshold {
		updates := make([]*models.RawTopicMessage,len(processor.topicUpdates))
		copy(updates,processor.topicUpdates)
		processor.topicUpdates = make([]*models.RawTopicMessage,0,processor.bulkUpdateThreshold)
		fmt.Println("Bulk Update")
		sortedUpdates := sortTopicUpdatesByName(updates)
		fetchRequest := make([]string,0,len(sortedUpdates))
		for name,_ := range sortedUpdates {
			fetchRequest = append(fetchRequest,name)
		}
		existingTopics := FindTopicsByName(fetchRequest)
		fmt.Printf("Number of Existing Topics: %d",len(existingTopics))
		processSortedTopics(existingTopics,sortedUpdates)
	}
}

func sortTopicUpdatesByName(topicUpdates []*models.RawTopicMessage) map[string][]*models.RawTopicMessage {
	sortedTopics := make(map[string][]*models.RawTopicMessage)
	for _,topic := range topicUpdates {
		topicArray, entryExists := sortedTopics[topic.Name]
		if entryExists {
			topicArray = append(topicArray,topic)
		} else {
			sortedTopics[topic.Name] = make([]*models.RawTopicMessage,0,len(topicUpdates))
			sortedTopics[topic.Name] = append(sortedTopics[topic.Name],topic)
		}

	}
	return sortedTopics
}

func processSortedTopics(existingTopics map[string]*models.Topic, sortedTopics map[string][]*models.RawTopicMessage) {
	resultingTopicUpdates := make([]*models.Topic, 0, len(sortedTopics))
	for name, topicArray := range sortedTopics {
		var resultingTopic *models.Topic
		existingTopic, _ := existingTopics[name]
		resultingTopic = existingTopic
		for _, topic := range topicArray {
			resultingTopic = updateTopicInformation(resultingTopic,topic)
		}
		resultingTopicUpdates = append(resultingTopicUpdates,resultingTopic)
	}

	StoreTopics(resultingTopicUpdates)
}

func updateTopicInformation(existingTopic *models.Topic, newTopic *models.RawTopicMessage) *models.Topic {
	var resultingTopic *models.Topic

	if existingTopic == nil {
		resultingTopic = models.NewTopic(newTopic.Name, newTopic.Payload)
		resultingTopic.CalculateUpdateBehavior(0)
	} else {
		resultingTopic = existingTopic
		newUpdateInterval := int(newTopic.ArrivalTime.Sub(resultingTopic.LastUpdateTimeStamp).Seconds())
		resultingTopic.LastPayload = newTopic.Payload
		resultingTopic.LastUpdateTimeStamp = newTopic.ArrivalTime
		resultingTopic.CalculateUpdateBehavior(newUpdateInterval)
	}
	return resultingTopic
}


