package common

import (
	"sync/atomic"
	"sync"
	"github.com/tkrex/IDS/common/models"
	"fmt"
)

type TopicProcessor struct {

	// Flag to indicate that the consumer state
	// 0 == Running
	// 1 == Closed
	state                int64
	processorStarted     sync.WaitGroup
	processorStopped     sync.WaitGroup
	incomingTopicChannel chan *models.Topic
	topicPersistenceManager InformationPersistenceManager
}

func NewTopicProcessor(persistenceManager InformationPersistenceManager, incomingTopicChannel chan *models.Topic ) *TopicProcessor {
	processor := new(TopicProcessor)
	processor.processorStarted.Add(1)
	processor.processorStopped.Add(1)
	processor.topicPersistenceManager = persistenceManager
	processor.incomingTopicChannel = incomingTopicChannel
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

	topic, ok :=  <-processor.incomingTopicChannel
	if topic != nil {
		var resultingTopic models.Topic
		var newUpdateInterval int
		if existingTopic, found := processor.topicPersistenceManager.TopicWithName(topic.Name); found {
			resultingTopic = existingTopic
			newUpdateInterval = int(topic.LastUpdateTimeStamp.Sub(resultingTopic.LastUpdateTimeStamp).Seconds())
			resultingTopic.LastPayload = topic.LastPayload
			resultingTopic.LastUpdateTimeStamp = topic.LastUpdateTimeStamp
		} else {
			resultingTopic = *topic
			newUpdateInterval = 0
		}

		resultingTopic.CalculateUpdateBehavior(newUpdateInterval)
		processor.topicPersistenceManager.StoreTopic(resultingTopic)
	}
	return ok
}