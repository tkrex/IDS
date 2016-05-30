package forwarding

import (
	"sync"
	"github.com/tkrex/IDS/common/models"
	"encoding/json"
	"fmt"
	"time"
	"github.com/tkrex/IDS/daemon/persistence"
	"github.com/tkrex/IDS/common/controlling"
	"github.com/tkrex/IDS/common/publishing"
)

type DomainInformationForwarder struct {
	forwarderStarted     sync.WaitGroup
	forwarderStopped     sync.WaitGroup
	forwardSignalChannel chan int
	lastForwardTimestamp           time.Time
}

const (
	ForwardInterval = 5 * time.Minute
	ForwardTopic = "DomainInformation"
)

func NewDomainInformationForwarder(forwardSignalChannel chan int) *DomainInformationForwarder {
	forwarder := new(DomainInformationForwarder)
	forwarder.forwardSignalChannel = forwardSignalChannel
	forwarder.lastForwardTimestamp = time.Now()
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

func (forwarder *DomainInformationForwarder) listenOnForwardSignal() {
	for {
		shouldForward, open := <- forwarder.forwardSignalChannel
		if !open {
			break
		}
		if shouldForward == 1 {
			go forwarder.forwardAllDomainInformation()
		}
	}
}

func (forwarder *DomainInformationForwarder) startForwardTicker() {
	forwardTicker := time.NewTicker(ForwardInterval)
	for _ = range forwardTicker.C {
		fmt.Println("Forward Ticker tick")
		forwarder.triggerForwarding()
	}
}

func (forwarder *DomainInformationForwarder) triggerForwarding() {
	if time.Now().Sub(forwarder.lastForwardTimestamp) > ForwardInterval {
		forwarder.forwardAllDomainInformation()
	}
}

func (forwarder *DomainInformationForwarder) forwardAllDomainInformation() {
	defer func() { forwarder.lastForwardTimestamp = time.Now()}()

	fmt.Println("Forwarding All Domain Information")
	dbDelegate, _ := persistence.NewDaemonDatabaseWorker()
	if dbDelegate == nil {
		return
	}
	domains, _ := dbDelegate.FindAllDomains()
	dbDelegate.Close()

	for _, domain := range domains {
		forwarder.forwardDomainInformation(domain)
	}

}

func (forwarder *DomainInformationForwarder) calculateForwardPriority(domainInformation *models.DomainInformationMessage) {
	priority := 0
	for _,topic := range domainInformation.Topics {
		if  topic.FirstUpdateTimeStamp.After(forwarder.lastForwardTimestamp) {
			priority++
		}
	}
	domainInformation.ForwardPriority = priority
}

func (forwarder *DomainInformationForwarder) forwardDomainInformation(domain *models.RealWorldDomain) {
	dbDelagte, err := persistence.NewDaemonDatabaseWorker()
	if err != nil {
		fmt.Println(err)
	}

	defer dbDelagte.Close()

	domainInformation := dbDelagte.FindDomainInformationByDomainName(domain.Name)


	if domainInformation == nil {
		fmt.Printf("\n No topics for domain %s found", domain.Name)
		return
	}
	fmt.Printf("\n Forwarding %d Topics for the domain %s", len(domainInformation.Topics), domainInformation.RealWorldDomain.Name)



	if len(domainInformation.Topics) == 0 {
		dbDelagte.RemoveDomain(domain)
		return
	}

	forwarder.calculateForwardPriority(domainInformation)

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
	defer controlMessagesDBDelagte.Close()
	if domainController := controlMessagesDBDelagte.FindDomainControllerForDomain(domain.Name); domainController != nil {
		serverAddress = domainController.IpAddress
		fmt.Println("Sending information to  DomainController: ",domainController.Domain.Name,domainController.IpAddress)
	} else if rootController := controlMessagesDBDelagte.FindDomainControllerForDomain("default"); rootController != nil {
		serverAddress = rootController.IpAddress
		fmt.Println("Sending information to Default DomainController: ", rootController.IpAddress)
	}

	if serverAddress == "" {
		fmt.Println("No Domain Controller found for forwarding")
		return
	}


	domainControllerPublisherConfig := models.NewMqttClientConfiguration(serverAddress,"1883","tcp", ForwardTopic, domainInformation.Broker.ID)
	domainControllerPublisher := publishing.NewMqttPublisher(domainControllerPublisherConfig,false)
	domainControllerPublisher.Publish(json)
	domainControllerPublisher.Close()

	//brokerPublisherConfig := models.NewMqttClientConfiguration("localhost","1883","ws", "IDSStatistics/"+domain.Name, domainInformation.Broker.ID)
	//brokerPublisher := publishing.NewMqttPublisher(brokerPublisherConfig,true)
	//brokerPublisher.Publish(json)
	//brokerPublisher.Close()
}