package layers

import (
	"sync"
	"github.com/tkrex/IDS/common/models"
	"github.com/tkrex/IDS/common/layers"
	"encoding/json"
	"fmt"
)

type DomainInformationForwarder struct {
	publisher          common.InformationPublisher

	forwarderStarted   sync.WaitGroup
	forwarderStopped   sync.WaitGroup

	forwardSignalChannel chan int
	databaseDelegate *DaemonDatabaseWorker

}

func NewDomainInformationForwarder(forwardSignalChannel chan int) *DomainInformationForwarder {
	forwarder := new(DomainInformationForwarder)
	forwarder.forwardSignalChannel = forwardSignalChannel
	forwarder.forwarderStarted.Add(1)
	forwarder.forwarderStopped.Add(1)
	go forwarder.run()
	forwarder.forwarderStarted.Wait()
	return forwarder
}

func (forwarder *DomainInformationForwarder) run() {
	config := models.NewMqttClientConfiguration("tcp://localhost:1883","domainController","publisher")
	forwarder.publisher = common.NewMqttPublisher(config)
	go forwarder.listenOnForwardSignal()
	forwarder.forwarderStarted.Done()
}

func (forwarder *DomainInformationForwarder) close() {
	forwarder.publisher.Close()
}

func (forwarder *DomainInformationForwarder) listenOnForwardSignal() {
	for {
		forwardStatus, open := <- forwarder.forwardSignalChannel
		if !open {
			break
		}
		if forwardStatus == 1 {
			go forwarder.forwardDomainInformation()
		}
	}
}

func (forwarder *DomainInformationForwarder) forwardDomainInformation() {
		dbDelegate,_ := NewDaemonDatabaseWorker()
		topics,_ := dbDelegate.FindAllTopics()
		broker,_ := dbDelegate.FindBroker()
		message := models.NewDomainInformationMessage(broker.RealWorldDomains[0],broker,topics)
		json, err := json.Marshal(message)
		if err != nil {
			fmt.Printf("Marshalling Error: %s",err)
			return
		}
		go forwarder.publisher.Publish(json)
}