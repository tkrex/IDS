package common

import (
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/tkrex/IDS/common/models"
	"sync"
	"sync/atomic"
	"fmt"
)

type MqttPublisher struct {
	state                 int64
	outgoingTopicsChannel chan *models.Topic
	client                mqtt.Client
	config *models.MqttClientConfiguration

	publisherStarted       sync.WaitGroup
	publisherStopped       sync.WaitGroup
}

func NewMqttPublisher(config *models.MqttClientConfiguration) *MqttPublisher {
	publisher := new(MqttPublisher)
	publisher.config = config
	publisher.publisherStarted.Add(1)
	publisher.publisherStopped.Add(1)
	go publisher.run()
	publisher.publisherStarted.Wait()
	return publisher
}

func (publisher *MqttPublisher) run() {

	opts := mqtt.NewClientOptions().AddBroker(publisher.config.BrokerAddress())
	opts.SetClientID(publisher.config.ClientID())
	publisher.client = mqtt.NewClient(opts)

	if token := publisher.client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	publisher.publisherStarted.Done()
}

func (publisher *MqttPublisher) PublishTopics(topics []models.Topic)  {
	for _,topic := range topics {
		if err := publisher.PublishTopic(topic); err != nil {
			fmt.Printf("Publish Error: %s",err)
		}
	}
}

func (publisher *MqttPublisher) PublishTopic(topic models.Topic) error {
	fmt.Println("Publishing:" + topic.Name)
	if token := publisher.client.Publish(publisher.config.Topic(), 2, false, topic.Name); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (publisher *MqttPublisher) Close() {
	atomic.StoreInt64(&publisher.state,1)
	publisher.client.Disconnect(10)
	fmt.Println("Publisher Disconnected")
}