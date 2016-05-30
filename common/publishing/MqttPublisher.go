package publishing

import (
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/tkrex/IDS/common/models"
	"sync"
	"sync/atomic"
	"fmt"
)

type MqttPublisher struct {
	state                 int64
	client                mqtt.Client
	config *models.MqttClientConfiguration
	retained bool

	publisherStarted       sync.WaitGroup
	publisherStopped       sync.WaitGroup
}

func NewMqttPublisher(config *models.MqttClientConfiguration, retained bool) *MqttPublisher {
	publisher := new(MqttPublisher)
	publisher.config = config
	publisher.retained = retained
	publisher.publisherStarted.Add(1)
	publisher.publisherStopped.Add(1)
	go publisher.run()
	publisher.publisherStarted.Wait()
	fmt.Println("Publisher started")
	return publisher
}

func (publisher *MqttPublisher) run() {

	opts := mqtt.NewClientOptions().AddBroker(publisher.config.BrokerAddress())
	opts.SetClientID(publisher.config.ClientID())
	opts.WillRetained = publisher.retained
	publisher.client = mqtt.NewClient(opts)

	if token := publisher.client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	publisher.publisherStarted.Done()
}



func (publisher *MqttPublisher) Publish(data []byte) error  {
	if token := publisher.client.Publish(publisher.config.Topic(), 2, false, data); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	fmt.Println("Published to Domain Controller")
	return nil
}

func (publisher *MqttPublisher) Close() {
	atomic.StoreInt64(&publisher.state,1)
	publisher.client.Disconnect(10)
	fmt.Println("Publisher Disconnected")
}