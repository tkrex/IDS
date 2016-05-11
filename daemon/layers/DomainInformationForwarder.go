package layers

import (
	"sync"
	"github.com/tkrex/IDS/common/models"
	"github.com/tkrex/IDS/common/layers"
	"encoding/json"
	"fmt"
	"time"
)

type DomainInformationForwarder struct {
	forwarderStarted     sync.WaitGroup
	forwarderStopped     sync.WaitGroup
	forwardSignalChannel chan int
	databaseDelegate     *DaemonDatabaseWorker
	lastForwardTimestamp           time.Time
}

const (
	ForwardInterval = 5 * time.Minute
	ForwardTopic = "DomainInformation"
)

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
	defer func() {
		forwarder.lastForwardTimestamp = time.Now()
	}()
	if time.Now().Sub(forwarder.lastForwardTimestamp) > ForwardInterval {
		forwarder.forwardAllDomainInformation()
	}
}

func (forwarder *DomainInformationForwarder) forwardAllDomainInformation() {
	fmt.Println("Forwarding All Domain Information")
	dbDelegate, _ := NewDaemonDatabaseWorker()
	if dbDelegate == nil {
		return
	}
	domains, _ := dbDelegate.FindAllDomains()
	dbDelegate.Close()

	for _, domain := range domains {
		go forwarder.forwardDomainInformation(domain)
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
	fmt.Println("Forward Priority: ",priority)
}
func (forwarder *DomainInformationForwarder) forwardDomainInformation(domain *models.RealWorldDomain) {
	dbDelagte, err := NewDaemonDatabaseWorker()
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
	if domainController, _ := dbDelagte.FindDomainControllerForDomain(domain.Name); domainController != nil {
		serverAddress = domainController.IpAddress
	} else if rootController, _ := dbDelagte.FindDomainControllerForDomain("rootController"); rootController != nil {
		serverAddress = rootController.IpAddress
	}

	if serverAddress == "" {
		fmt.Println("No Domain Controller found for forwarding")
		return
	}

	publisherConfig := models.NewMqttClientConfiguration(serverAddress, ForwardTopic, domainInformation.Broker.ID)
	publisher := common.NewMqttPublisher(publisherConfig)
	publisher.Publish(json)
	publisher.Close()
}