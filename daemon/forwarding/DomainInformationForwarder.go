package forwarding

import (
	"sync"
	"github.com/tkrex/IDS/common/models"
	"encoding/json"
	"fmt"
	"time"
	"github.com/tkrex/IDS/daemon/persistence"
	"github.com/tkrex/IDS/common/publishing"
	"github.com/tkrex/IDS/common/forwardRouting"
	"github.com/tkrex/IDS/daemon/configuration"
)

//Forwards DomainInformationMessages to the corresponding Domain Controller
type DomainInformationForwarder struct {
	forwarderStarted     sync.WaitGroup
	forwarderStopped     sync.WaitGroup
	forwardSignalChannel chan int
	lastForwardTimestamp           time.Time
	dbDelegate *persistence.DomainInformationStorage
	routingManager 	*forwardRouting.ForwardRoutingManager
}

const (
	ForwardInterval = 5 * time.Minute
)

func NewDomainInformationForwarder(forwardSignalChannel chan int) *DomainInformationForwarder {
	forwarder := new(DomainInformationForwarder)
	forwarder.forwardSignalChannel = forwardSignalChannel
	forwarder.lastForwardTimestamp = time.Now()
	forwarder.routingManager = forwardRouting.NewForwardRoutingManager(configuration.DaemonConfigurationManagerInstance().Config().RoutingManagementURL)
	forwarder.forwarderStarted.Add(1)
	forwarder.forwarderStopped.Add(1)
	go forwarder.run()
	forwarder.forwarderStarted.Wait()
	return forwarder
}

//Connects to database and starts listing on channels for signals
func (forwarder *DomainInformationForwarder) run() {
	dbDeleagte , error := persistence.NewDomainInformationStorage()
	if error != nil {
		fmt.Println("Cannot connect to Database. STopping Forwarder.")
		return
	}
	forwarder.dbDelegate = dbDeleagte
	go forwarder.listenOnForwardSignal()
	go forwarder.startForwardTicker()
	forwarder.forwarderStarted.Done()
}

//Starts listening on Forward Signals from Topic Processor
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

//Start Tickers which triggers the forwarding of DomainInformationMessages when expiring
func (forwarder *DomainInformationForwarder) startForwardTicker() {
	forwardTicker := time.NewTicker(ForwardInterval)
	for _ = range forwardTicker.C {
		fmt.Println("Forward Ticker tick")
		forwarder.forwardAllDomainInformation()
	}
}


//Initiate the forwarding process for each available Real World Domain
func (forwarder *DomainInformationForwarder) forwardAllDomainInformation() {
	defer func() { forwarder.lastForwardTimestamp = time.Now()}()

	fmt.Println("Forwarding All Domain Information")
	dbDelegate, _ := persistence.NewDomainInformationStorage()
	if dbDelegate == nil {
		return
	}
	domains, _ := dbDelegate.FindAllDomains()
	dbDelegate.Close()

	for _, domain := range domains {
		forwarder.forwardDomainInformation(domain)
	}

}

//Calculates how many topics were newly added since the last forwarding. The resulting value serves as indicator for the importance of the message.
// The higher the value the faster the message is forwarded to the Top Level Domain Controller
func (forwarder *DomainInformationForwarder) calculateForwardPriority(domainInformation *models.DomainInformationMessage) {
	priority := 0
	for _,topic := range domainInformation.Topics {
		if  topic.FirstUpdateTimeStamp.After(forwarder.lastForwardTimestamp) {
			priority++
		}
	}
	domainInformation.ForwardPriority = priority
}


//Creates an DomainInformationMessage for the specified Real World Domain, requests the address of the Domain Controller and forwards the message to this address.
func (forwarder *DomainInformationForwarder) forwardDomainInformation(domain *models.RealWorldDomain) {

	domainInformation := forwarder.dbDelegate.FindDomainInformationByDomainName(domain.Name)


	if domainInformation == nil {
		fmt.Printf("\n No topics for domain %s found", domain.Name)
		return
	}
	fmt.Printf("\n Forwarding %d Topics for the domain %s", len(domainInformation.Topics), domainInformation.RealWorldDomain.Name)



	if len(domainInformation.Topics) == 0 {
		forwarder.dbDelegate.RemoveDomain(domain)
		return
	}

	forwarder.calculateForwardPriority(domainInformation)

	json, err := json.Marshal(domainInformation)
	if err != nil {
		fmt.Printf("Marshalling Error: %s", err)
		return
	}

	domainController,err := forwarder.routingManager.DomainControllerForDomain(domain,false)
	if err != nil {
		fmt.Println("Forwarder: No Target Controller Found")
		return
	}

	domainControllerPublisherConfig := models.NewMqttClientConfiguration(domainController.BrokerAddress, domainInformation.Broker.ID)
	domainControllerPublisher := publishing.NewMqttPublisher(domainControllerPublisherConfig,false)
	error := domainControllerPublisher.Publish(json,domainInformation.Broker.ID)
	if error != nil {
		domainController,err := forwarder.routingManager.DomainControllerForDomain(domain,true)
		if err != nil {
			fmt.Println("Forwarder: No Target Controller Found")
			return
		}
		domainControllerPublisherConfig := models.NewMqttClientConfiguration(domainController.BrokerAddress, domainInformation.Broker.ID)
		domainControllerPublisher := publishing.NewMqttPublisher(domainControllerPublisherConfig,false)
		error := domainControllerPublisher.Publish(json,domainInformation.Broker.ID)
		if error != nil {
			fmt.Println(error)
			return
		}
	}
	domainControllerPublisher.Close()
}

