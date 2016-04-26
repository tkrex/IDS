package common

import (
	"sync/atomic"
	"sync"
	"github.com/tkrex/IDS/common/models"
	"fmt"
	"github.com/tkrex/IDS/common"
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
		if existingTopic, found := processor.topicPersistenceManager.TopicWithName(topic.Name); found {
			newUpdateInterval := topic.LastUpdateTimeStamp.Sub(existingTopic.LastUpdateTimeStamp).Seconds()
			calculateUpdateBehavior(&existingTopic.UpdateBehavior,int(newUpdateInterval))
			existingTopic.LastPayload = topic.LastPayload
			existingTopic.LastUpdateTimeStamp = topic.LastUpdateTimeStamp
			processor.topicPersistenceManager.StoreTopic(existingTopic)
		} else {
			processor.topicPersistenceManager.StoreTopic(*topic)
		}
	}
	return ok
}

func calculateUpdateBehavior(updateBehavior *models.UpdateBehavior, newUpdateInterval int) {
	if updateBehavior.NumberOfUpdates == 0 {
		updateBehavior.AverageUpdateIntervalInSeconds = newUpdateInterval
		updateBehavior.MaximumUpdateIntervalInSeconds = newUpdateInterval
		updateBehavior.MinimumUpdateIntervalInSeconds = newUpdateInterval
	} else {
		updateBehavior.MaximumUpdateIntervalInSeconds = common.Max(updateBehavior.MaximumUpdateIntervalInSeconds,newUpdateInterval)
		updateBehavior.MinimumUpdateIntervalInSeconds = common.Min(updateBehavior.MinimumUpdateIntervalInSeconds,newUpdateInterval)
		updateBehavior.AverageUpdateIntervalInSeconds = (updateBehavior.AverageUpdateIntervalInSeconds * updateBehavior.NumberOfUpdates + newUpdateInterval) / (updateBehavior.NumberOfUpdates +1)
	}
	updateBehavior.NumberOfUpdates++
}