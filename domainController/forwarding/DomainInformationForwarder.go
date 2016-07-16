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
)

type DomainInformationForwarder struct {
	forwarderStarted       sync.WaitGroup
	forwarderStopped       sync.WaitGroup

	forwardSignalChannel   chan *models.RealWorldDomain
	routingManager         *routing.RoutingManager
	forwardPriorityCounter map[string]int
}

const (
	DomainForwardInterval = 1 * time.Hour
	ForwardThreshold = 10
)

func NewDomainInformationForwarder(forwardSignalChannel chan *models.ForwardMessage) *DomainInformationForwarder {
	forwarder := new(DomainInformationForwarder)
	forwarder.forwardSignalChannel = forwardSignalChannel
	forwarder.routingManager = routing.NewRoutingManager(configuration.DomainControllerConfigurationManagerInstance().Config().ScalingInterfaceAddress)
	forwarder.forwardPriorityCounter = make(map[string]int)
	forwarder.forwarderStarted.Add(1)
	forwarder.forwarderStopped.Add(1)
	go forwarder.run()
	forwarder.forwarderStarted.Wait()
	return forwarder
}

func (forwarder *DomainInformationForwarder) run() {
	go forwarder.listenOnForwardSignal()
	go forwarder.startForwardTicker()
	forwarder.forwarderStarted.Done()
}

func (forwarder *DomainInformationForwarder) close() {
}

func (forwarder *DomainInformationForwarder) listenOnForwardSignal() {
	for {
		forwardMessage, open := <-forwarder.forwardSignalChannel
		if !open {
			break
		}
		if forwardMessage != nil {
			go forwarder.processForwardMessage(forwardMessage)
		}
	}
}

func (forwarder *DomainInformationForwarder) startForwardTicker() {
	forwardTicker := time.NewTicker(DomainForwardInterval)
	for _ = range forwardTicker.C {
		forwarder.checkDomainsForForwarding()
	}
}

func (forwarder *DomainInformationForwarder) processForwardMessage(forwardMessage *models.ForwardMessage) {
	config := configuration.DomainControllerConfigurationManagerInstance().Config()
	parentDomain := config.ParentDomain

	domain := forwardMessage.Domain
	priority := forwardMessage.Priority

	if domain.IsSubDomainOf(parentDomain) {
		forwarder.forwardPriorityCounter[domain.Name] += priority
		if forwarder.forwardPriorityCounter[domain.Name] >= ForwardThreshold {
			forwarder.forwardDomainInformation(domain)

			return
		}
	} else {
		forwarder.forwardDomainInformation(domain)
	}
}

func (forwarder *DomainInformationForwarder) checkDomainsForForwarding() {
	dbDelagte, _ := persistence.NewDomainControllerDatabaseWorker()
	defer dbDelagte.Close()
	domains, _ := dbDelagte.FindAllDomains()
	for _, domain := range domains {
		forwarder.forwardDomainInformation(domain)
	}
}

func (forwarder *DomainInformationForwarder) forwardDomainInformation(domain *models.RealWorldDomain) {
	forwarder.forwardPriorityCounter[domain.Name] = 0

	domainInformationDelegate, _ := persistence.NewDomainControllerDatabaseWorker()
	defer domainInformationDelegate.Close()

	domainInformation, err := domainInformationDelegate.FindDomainInformationByDomainName(domain.Name)

	if err != nil {
		fmt.Println(err)
		return
	}

	if len(domainInformation) == 0 {
		domainInformationDelegate.RemoveDomain(domain)
		return
	}

	configManager := configuration.DomainControllerConfigurationManagerInstance()
	targetDomain := new(models.RealWorldDomain)

	parentDomain := configManager.Config().ParentDomain
	ownDomain := configManager.Config().OwnDomain
	if ownDomain.Name == "default" {
		targetDomain = domain
	} else {
		targetDomain = parentDomain
	}

	domainController, err := forwarder.routingManager.DomainControllerForDomain(targetDomain, false)
	if err != nil {
		fmt.Println("Forwarder: No Target Controller Found")
		return
	}

	domainControllerPublisherConfig := models.NewMqttClientConfiguration(domainController.BrokerAddress, configManager.Config().DomainControllerID)
	domainControllerPublisher := publishing.NewMqttPublisher(domainControllerPublisherConfig, false)

	for _, information := range domainInformation {
		json, err := json.Marshal(information)
		if err != nil {
			fmt.Printf("Marshalling Error: %s", err)
			return
		}
		error := domainControllerPublisher.Publish(json, information.Broker.ID)
		if error != nil {
			domainController, err := forwarder.routingManager.DomainControllerForDomain(targetDomain, true)
			if err != nil {
				fmt.Println("Forwarder: No Target Controller Found")
				return
			}
			domainControllerPublisherConfig := models.NewMqttClientConfiguration(domainController.BrokerAddress, information.Broker.ID)
			domainControllerPublisher := publishing.NewMqttPublisher(domainControllerPublisherConfig, false)
			error := domainControllerPublisher.Publish(json, information.Broker.ID)
			if error != nil {
				fmt.Println(error)
				return
			}
		}
		domainControllerPublisher.Close()
	}
}