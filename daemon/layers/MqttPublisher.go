package layers

import (
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/tkrex/IDS/common/models"
	"sync"
	"sync/atomic"
	"time"
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
	for closed := atomic.LoadInt64(&publisher.state) == 1; !closed; closed = atomic.LoadInt64(&publisher.state) == 1 {
		time.Sleep(time.Second * 5)
		//Publish
		fmt.Println("Publish Message")
		if token := publisher.client.Publish(publisher.publishedTopic, 2, false, "Test"); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
			os.Exit(1)
		}

	}
	publisher.publisherStopped.Done()
}



func (publisher *MqttPublisher) Close() {
	atomic.StoreInt64(&publisher.state,1)
	publisher.publisherStopped.Wait()
	publisher.client.Disconnect(10)
	fmt.Println("Publisher Disconnected")
}

