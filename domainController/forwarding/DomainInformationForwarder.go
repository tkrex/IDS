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
	forwarderStarted     sync.WaitGroup
	forwarderStopped     sync.WaitGroup

	forwardSignalChannel chan *models.RealWorldDomain
	updateFlags          map[string]bool

	routingManager *routing.RoutingManager
}

const (
	ForwardInterval = 1 * time.Minute
)

func NewDomainInformationForwarder(forwardSignalChannel chan *models.RealWorldDomain) *DomainInformationForwarder {
	forwarder := new(DomainInformationForwarder)
	forwarder.forwardSignalChannel = forwardSignalChannel
	forwarder.routingManager = routing.NewRoutingManager()
	forwarder.updateFlags = make(map[string]bool)
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
		domain, open := <-forwarder.forwardSignalChannel
		if !open {
			break
		}
		if domain != nil {
			go forwarder.forwardDomainInformation(domain)
		}
	}
}

func (forwarder *DomainInformationForwarder) startForwardTicker() {
	forwardTicker := time.NewTicker(ForwardInterval)
	for _ = range forwardTicker.C {
		forwarder.checkDomainsForForwarding()
	}
}

func (forwarder *DomainInformationForwarder) checkDomainsForForwarding() {
	dbDelagte, _ := persistence.NewDomainControllerDatabaseWorker()
	defer dbDelagte.Close()
	domains, _ := dbDelagte.FindAllDomains()
	for _, domain := range domains {
		if updateFlag := forwarder.updateFlags[domain.Name]; !updateFlag {
			go forwarder.forwardDomainInformation(domain)
		}
		forwarder.updateFlags[domain.Name] = false
	}

}

func (forwarder *DomainInformationForwarder) forwardDomainInformation(domain *models.RealWorldDomain) {
	domainInformationDelegate, _ := persistence.NewDomainControllerDatabaseWorker()
	defer domainInformationDelegate.Close()

	domainInformation, err := domainInformationDelegate.FindDomainInformationByDomainName(domain.Name)

	if err != nil {
		fmt.Println(err)
		return
	}

	if len(domainInformation) == 0 {
		domainInformationDelegate.RemoveDomain(domain)
		delete(forwarder.updateFlags, domain.Name)
		return
	}


	configManager := configuration.NewDomainControllerConfigurationManager()
	config,_ := configManager.DomainControllerConfig()
	targetDomain := new(models.RealWorldDomain)

	parentDomain := config.ParentDomain
	if !domain.IsSubDomainOf(parentDomain) {
		targetDomain = domain
	} else {
		targetDomain = parentDomain
	}


	domainController,err := forwarder.routingManager.DomainControllerForDomain(targetDomain,false)
	if err != nil {
		fmt.Println("Forwarder: No Target Controller Found")
		return
	}

	domainControllerPublisherConfig := models.NewMqttClientConfiguration(domainController.BrokerAddress, config.DomainControllerID)
	domainControllerPublisher := publishing.NewMqttPublisher(domainControllerPublisherConfig,false)


	for _, information := range domainInformation {
		json, err := json.Marshal(information)
		if err != nil {
			fmt.Printf("Marshalling Error: %s", err)
			return
		}
		error := domainControllerPublisher.Publish(json, information.Broker.ID)
		if error != nil {
			domainController,err := forwarder.routingManager.DomainControllerForDomain(targetDomain,true)
			if err != nil {
				fmt.Println("Forwarder: No Target Controller Found")
				return
			}
			domainControllerPublisherConfig := models.NewMqttClientConfiguration(domainController.BrokerAddress, information.Broker.ID)
			domainControllerPublisher := publishing.NewMqttPublisher(domainControllerPublisherConfig,false)
			error := domainControllerPublisher.Publish(json,information.Broker.ID)
			if error != nil {
				fmt.Println(error)
				return
			}
		}
		domainControllerPublisher.Close()
	}
}