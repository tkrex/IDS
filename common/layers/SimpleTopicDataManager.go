package common

import (
	"sync/atomic"
	"sync"
	"github.com/tkrex/IDS/common/models"
	"fmt"
)

type SimpleTopicDataManager struct {

	// Flag to indicate that the consumer state
	// 0 == Running
	// 1 == Closed
	state                int64
	consumerStarted sync.WaitGroup
	consumerStopped sync.WaitGroup
	incomingTopicChannel chan *models.Topic
	topicsCollection     map[string]*models.Topic
}

func  NewDataManager(producer InformationProducer) *SimpleTopicDataManager {
	manager := new(SimpleTopicDataManager)
	manager.incomingTopicChannel = producer.InformationChannel()
	manager.topicsCollection = make(map[string]*models.Topic)
	manager.consumerStarted.Add(1)
	manager.consumerStopped.Add(1)

	go manager.consumer()
	manager.consumerStarted.Wait()
	fmt.Println("Consumer Created")
	return manager
}


func (consumer *SimpleTopicDataManager) State() int64 {
	return atomic.LoadInt64(&consumer.state)
}


func (consumer* SimpleTopicDataManager)  Close() {
	fmt.Println("Closing Consumer")
	atomic.StoreInt64(&consumer.state,1)
	consumer.consumerStopped.Wait()
	fmt.Println("Consumer Closed")
}

func  (consumer *SimpleTopicDataManager) consumer() {
	consumer.consumerStarted.Done()

	for closed := atomic.LoadInt64(&consumer.state) == 1; !closed; closed = atomic.LoadInt64(&consumer.state) == 1 {
		  topic, ok :=  <-consumer.incomingTopicChannel
		if topic != nil {
			consumer.Store(topic)
		}
		if !ok {
			fmt.Println("IncomingTopicsChannel closed")
		}
	}
	consumer.consumerStopped.Done()
}

func (consumer *SimpleTopicDataManager) Store(topic *models.Topic) {
	fmt.Println(topic.Name)
	if existingTopic, ok := consumer.topicsCollection[topic.Name]; ok {
		newUpdateInterval := topic.LastUpdateTimeStamp.Sub(existingTopic.LastUpdateTimeStamp).Seconds()
		existingTopic.UpdateInterval = int(newUpdateInterval)
		existingTopic.LastPayload = topic.LastPayload
		existingTopic.LastUpdateTimeStamp = topic.LastUpdateTimeStamp
		existingTopic.NumberOfUpdates++
		consumer.topicsCollection[topic.Name] = existingTopic
	} else {
		consumer.topicsCollection[topic.Name] = topic

	}
}
