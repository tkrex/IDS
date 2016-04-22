package common

import (
	"sync/atomic"
	"sync"
	"github.com/tkrex/IDS/common/models"
	"fmt"
	"time"
	"github.com/tkrex/IDS/daemon/layers"
)

type TopicProcessor struct {

	// Flag to indicate that the consumer state
	// 0 == Running
	// 1 == Closed
	state                int64
	processorStarted     sync.WaitGroup
	processorStopped     sync.WaitGroup
	incomingTopicChannel chan *models.Topic
	topicsCollection     map[string]*models.Topic
	topicPublisher 	     InformationPublisher
	publishTicker	*time.Ticker


}

func NewTopicProcessor(producer InformationProducer) *TopicProcessor {
	processor := new(TopicProcessor)
	processor.incomingTopicChannel = producer.InformationChannel()
	processor.topicsCollection = make(map[string]*models.Topic)
	processor.processorStarted.Add(1)
	processor.processorStopped.Add(1)

	go processor.run()
	processor.processorStarted.Wait()
	fmt.Println("Consumer Created")
	return processor
}


func (processor *TopicProcessor) State() int64 {
	return atomic.LoadInt64(&processor.state)
}


func (processor *TopicProcessor)  Close() {
	fmt.Println("Closing Processor")
	processor.publishTicker.Stop()
	processor.topicPublisher.Close()
	atomic.StoreInt64(&processor.state,1)
	processor.processorStopped.Wait()
	fmt.Println("Processor Closed")
}

func  (processor *TopicProcessor) run() {

	processor.topicPublisher = layers.NewMqttPublisher("tcp://localhost:1883", "burst")
	processor.processorStarted.Done()

	processor.forwardTopics()
	for closed := atomic.LoadInt64(&processor.state) == 1; !closed; closed = atomic.LoadInt64(&processor.state) == 1 {
		processor.storeTopic()
	}
	processor.processorStopped.Done()
}

func (processor *TopicProcessor) forwardTopics() {
	processor.publishTicker = time.NewTicker(time.Second * 10)
	go func() {
		for _ = range processor.publishTicker.C {
			fmt.Println("Tick")
			topics := make(map[string]*models.Topic)
			for k, v := range processor.topicsCollection {
				topics[k] = v
			}
			go processor.topicPublisher.Publish(topics)
		}
	}()
}
func (processor *TopicProcessor) storeTopic() {

	topic, ok :=  <-processor.incomingTopicChannel
	if topic != nil {
		if existingTopic, ok := processor.topicsCollection[topic.Name]; ok {
			newUpdateInterval := topic.LastUpdateTimeStamp.Sub(existingTopic.LastUpdateTimeStamp).Seconds()
			existingTopic.UpdateInterval = int(newUpdateInterval)
			existingTopic.LastPayload = topic.LastPayload
			existingTopic.LastUpdateTimeStamp = topic.LastUpdateTimeStamp
			existingTopic.NumberOfUpdates++
			processor.topicsCollection[topic.Name] = existingTopic
		} else {
			processor.topicsCollection[topic.Name] = topic
		}
	}
	if !ok {
		fmt.Println("IncomingTopicsChannel closed")
	}
}
