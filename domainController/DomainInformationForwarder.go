package domainController

import (
	"sync"
	"github.com/tkrex/IDS/common/models"
	"github.com/tkrex/IDS/common/layers"
	"encoding/json"
	"fmt"
	"time"
)

type DomainInformationForwarder struct {
	publisher            common.InformationPublisher

	forwarderStarted     sync.WaitGroup
	forwarderStopped     sync.WaitGroup

	forwardSignalChannel chan *models.RealWorldDomain
	databaseDelegate     *DomainControllerDatabaseWorker

	updateFlags          map[string]bool
}

const (
	ForwardInterval = 1 * time.Minute
	ForwardTopic = "DomainInformation"
)

func NewDomainInformationForwarder(forwardSignalChannel chan *models.RealWorldDomain) *DomainInformationForwarder {
	forwarder := new(DomainInformationForwarder)
	forwarder.forwardSignalChannel = forwardSignalChannel
	forwarder.updateFlags = make(map[string]bool)
	forwarder.forwarderStarted.Add(1)
	forwarder.forwarderStopped.Add(1)
	go forwarder.run()
	forwarder.forwarderStarted.Wait()
	return forwarder
}

func (forwarder *DomainInformationForwarder) run() {
	config := models.NewMqttClientConfiguration("tcp://localhost:1883", "domainController", "publisher")
	forwarder.publisher = common.NewMqttPublisher(config)
	go forwarder.listenOnForwardSignal()
	go forwarder.startForwardTicker()
	forwarder.forwarderStarted.Done()
}

func (forwarder *DomainInformationForwarder) close() {
	forwarder.publisher.Close()
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
	dbDelagte, _ := NewDomainControllerDatabaseWorker()
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
	domainInformationDelegate, _ := NewDomainControllerDatabaseWorker()
	defer domainInformationDelegate.Close()

	domainInformation, err := domainInformationDelegate.FindDomainInformationByDomainName(domain.Name)

	if err != nil {
		fmt.Println(err)
		return
	}

	json, err := json.Marshal(domainInformation)
	if err != nil {
		fmt.Printf("Marshalling Error: %s", err)
		return
	}
	serverAddress := ""

	controlMessagesDBDelagte, err := common.NewControlMessageDBDelegate()
	if err != nil {
		fmt.Println(err)
		return
	}

	if domainController, _ := controlMessagesDBDelagte.FindDomainControllerForDomain(domain.Name); domainController != nil {
		serverAddress = domainController.IpAddress
	} else if rootController, _ := controlMessagesDBDelagte.FindDomainControllerForDomain("rootController"); rootController != nil {
		serverAddress = rootController.IpAddress
	}

	if serverAddress == "" {
		fmt.Println("No Domain Controller found for forwarding")
		return
	}

	//TODO: Come up with DomainConroller ID
	publisherConfig := models.NewMqttClientConfiguration(serverAddress, ForwardTopic, "DomainCOntrollerID")
	publisher := common.NewMqttPublisher(publisherConfig)
	publisher.Publish(json)
	publisher.Close()
}