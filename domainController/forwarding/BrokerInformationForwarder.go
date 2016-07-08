package forwarding

import (
	"sync"
	"encoding/json"
	"fmt"
	"time"
	"github.com/tkrex/IDS/common/models"
	"github.com/tkrex/IDS/domainController/persistence"
	"github.com/tkrex/IDS/common/publishing"
	"github.com/tkrex/IDS/common/routing"
	"github.com/tkrex/IDS/domainController/configuration"
	"github.com/tkrex/IDS/domainController"
)

type BrokerInformationForwarder struct {
	forwarderStarted     sync.WaitGroup
	forwarderStopped     sync.WaitGroup

}

const (
	ForwardInterval = 1 * time.Minute
)

func NewBrokerInformationForwarder() *BrokerInformationForwarder {
	forwarder := new(BrokerInformationForwarder)
	forwarder.forwarderStarted.Add(1)
	forwarder.forwarderStopped.Add(1)
	go forwarder.run()
	forwarder.forwarderStarted.Wait()
	return forwarder
}

func (forwarder *BrokerInformationForwarder) run() {
	go forwarder.startForwardTicker()
	forwarder.forwarderStarted.Done()
}

func (forwarder *BrokerInformationForwarder) close() {
}


func (forwarder *BrokerInformationForwarder) startForwardTicker() {
	forwardTicker := time.NewTicker(ForwardInterval)
	for _ = range forwardTicker.C {
		forwarder.forwardBrokerInformation()
	}
}


func (forwarder *BrokerInformationForwarder) forwardBrokerInformation() {
	dbManager, _ := persistence.NewDomainControllerDatabaseWorker()
	defer dbManager.Close()

	brokers, err := dbManager.FindAllBrokers()

	if err != nil {
		fmt.Println(err)
		return
	}

	configManager := configuration.NewDomainControllerConfigurationManager()
	config,_ := configManager.DomainControllerConfig()

	jsonArray,_ := json.Marshal(brokers)

	brokerPublisherConfig := models.NewMqttClientConfiguration(config.GatewayBrokerAddress, config.DomainControllerID)
	brokerPublisher := publishing.NewMqttPublisher(brokerPublisherConfig,false)
	brokerPublisher.Publish(jsonArray,"brokers")
}