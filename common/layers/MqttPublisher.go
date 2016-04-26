package common

import (
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/tkrex/IDS/common/models"
	"sync"
	"sync/atomic"
	"fmt"
	"encoding/json"
)

type MqttPublisher struct {
	state                 int64
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

func (publisher *MqttPublisher) PublishTopics(topics []models.Topic) error  {

	broker := models.NewBroker(3,"127.0.0.1","krex.com")
	message := models.NewTopicInformationMessage(broker,topics)

	json, err := json.Marshal(message)
	if err != nil {
		fmt.Printf("Marshalling Error: %s",err)
		return  err
	}
	if token := publisher.client.Publish(publisher.config.Topic(), 2, false, json); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (publisher *MqttPublisher) Close() {
	atomic.StoreInt64(&publisher.state,1)
	publisher.client.Disconnect(10)
	fmt.Println("Publisher Disconnected")
}