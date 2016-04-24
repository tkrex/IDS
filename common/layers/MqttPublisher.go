package common

import (
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/tkrex/IDS/common/models"
	"sync"
	"sync/atomic"
	"fmt"
	"os"
)

type MqttPublisher struct {
	state                 int64
	outgoingTopicsChannel chan *models.Topic
	client                mqtt.Client

	publisherStarted       sync.WaitGroup
	publisherStopped       sync.WaitGroup

	publishedTopic          string
	brokerAddress		string
}

func NewMqttPublisher(brokerAddress string, publishedTopic string) *MqttPublisher {
	publisher := new(MqttPublisher)
	publisher.publishedTopic = publishedTopic
	publisher.brokerAddress = brokerAddress
	publisher.publisherStarted.Add(1)
	publisher.publisherStopped.Add(1)
	go publisher.run()
	publisher.publisherStarted.Wait()
	return publisher
}

func (publisher *MqttPublisher) run() {

	opts := mqtt.NewClientOptions().AddBroker(publisher.brokerAddress)
	opts.SetClientID("publisher")
	publisher.client = mqtt.NewClient(opts)

	if token := publisher.client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	publisher.publisherStarted.Done()
}


func (publisher *MqttPublisher) PublishTopics(topics []*models.Topic) {
	for _,topic := range topics {
		fmt.Println("Publishing:" + topic.Name)
		if token := publisher.client.Publish(publisher.publishedTopic, 2, false, topic.Name); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
			os.Exit(1)
		}
	}

}

func (publisher *MqttPublisher) PublishTopic(topic *models.Topic) {
	fmt.Println("Publishing:" + topic.Name)
	if token := publisher.client.Publish(publisher.publishedTopic, 2, false, topic.Name); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}


func (publisher *MqttPublisher) Close() {
	atomic.StoreInt64(&publisher.state,1)
	publisher.client.Disconnect(10)
	fmt.Println("Publisher Disconnected")
}

