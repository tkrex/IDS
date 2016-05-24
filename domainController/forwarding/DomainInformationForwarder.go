package forwarding

import (
	"sync"
	"encoding/json"
	"fmt"
	"time"
	"github.com/tkrex/IDS/common/models"
	"github.com/tkrex/IDS/common/controlling"
	"github.com/tkrex/IDS/domainController/persistence"
	"github.com/tkrex/IDS/common/publishing"
)

type DomainInformationForwarder struct {
	forwarderStarted     sync.WaitGroup
	forwarderStopped     sync.WaitGroup

	forwardSignalChannel chan *models.RealWorldDomain
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
		delete(forwarder.updateFlags,domain.Name)
		return
	}

	json, err := json.Marshal(domainInformation)
	if err != nil {
		fmt.Printf("Marshalling Error: %s", err)
		return
	}
	serverAddress := ""

	controlMessagesDBDelagte, err := controlling.NewControlMessageDBDelegate()
	if err != nil {
		fmt.Println(err)
		return
	}


	//TODO: Get ParentDomain From ENV
	parentDomain := "default"

	if domainController := controlMessagesDBDelagte.FindDomainControllerForDomain(parentDomain); domainController != nil {
		serverAddress = domainController.IpAddress
	} else if rootController := controlMessagesDBDelagte.FindDomainControllerForDomain("default"); rootController != nil {
		serverAddress = rootController.IpAddress
	}

	if serverAddress == "" {
		fmt.Println("No Domain Controller found for forwarding")
		return
	}

	//TODO: Come up with DomainConroller ID
	publisherConfig := models.NewMqttClientConfiguration(serverAddress, ForwardTopic, "DomainControllerID")
	publisher :=  publishing.NewMqttPublisher(publisherConfig)
	publisher.Publish(json)
	publisher.Close()
}