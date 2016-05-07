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
	databaseDelegate *DaemonDatabaseWorker
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
	dbDelegate, error := NewDaemonDatabaseWorker()
	if error !=nil {
		fmt.Println("Stopping Forwarder: Cannot Connect to DB")
		return
	}
	forwarder.databaseDelegate = dbDelegate
	forwarder.forwarderStarted.Done()
	go func() {
		defer forwarder.databaseDelegate.Close()
		for _ = range forwarder.forwardTicker.C {
			topics,_ := dbDelegate.FindAllTopics()
			broker := models.NewBroker("127.0.0.1","krex.com")
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