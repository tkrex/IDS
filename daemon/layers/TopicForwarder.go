package layers

import (
	"time"
	"sync"
	"github.com/tkrex/IDS/common/models"
	"github.com/tkrex/IDS/common/layers"
	"encoding/json"
	"fmt"
)

type TopicForwarder struct {
	publisher          common.InformationPublisher
	forwardInterval    time.Duration
	forwardTicker      *time.Ticker

	forwarderStarted   sync.WaitGroup
	forwarderStopped   sync.WaitGroup
}

func NewTopicForwarder(forwardInterval time.Duration) *TopicForwarder {
	forwarder := new(TopicForwarder)
	forwarder.forwardInterval = forwardInterval
	forwarder.forwarderStarted.Add(1)
	forwarder.forwarderStopped.Add(1)
	go forwarder.run()
	forwarder.forwarderStarted.Wait()
	return forwarder
}

func (forwarder *TopicForwarder) run() {
	forwarder.forwardTicker = time.NewTicker(forwarder.forwardInterval)
	config := models.NewMqttClientConfiguration("tcp://localhost:1883","domainController","publisher")
	forwarder.publisher = common.NewMqttPublisher(config)
	forwarder.forwarderStarted.Done()
	go func() {
		for _ = range forwarder.forwardTicker.C {
			topics := FindAllTopics()
			broker := models.NewBroker(3,"127.0.0.1","krex.com")
			domain := models.NewRealWorldDomain("testDomain")
			message := models.NewDomainInformationMessage(domain,broker,topics)
			json, err := json.Marshal(message)
			if err != nil {
				fmt.Printf("Marshalling Error: %s",err)
				return
			}
			go forwarder.publisher.Publish(json)
		}
	}()
}