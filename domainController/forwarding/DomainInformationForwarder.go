package forwarding

import (
	"sync"
	"encoding/json"
	"fmt"
	"time"
	"github.com/tkrex/IDS/common/models"
	"github.com/tkrex/IDS/domainController/persistence"
	"github.com/tkrex/IDS/common/publishing"
	"github.com/tkrex/IDS/domainController/configuration"
	"github.com/tkrex/IDS/common/forwardRouting"
	"os"
)

//Mangages Forwarding of DomainInformationMessages to Parent Domain Controller
type DomainInformationForwarder struct {
	forwarderStarted       sync.WaitGroup
	forwarderStopped       sync.WaitGroup

	forwardSignalChannel   chan *models.ForwardMessage
	routingInformationChannel chan *models.DomainController
	routingManager         *forwardRouting.ForwardRoutingManager
	forwardPriorityCounter map[string]int

	informationStorageDelegate *persistence.DomainInformationStorage
}

const (
	//Timeout for forwarding Domain Information
	DomainForwardInterval = 1 * time.Hour
	//If sum of forward priorities of incoming DomainInformationMessages exceeds this threshold all DomainInformationMessages are forwarded
	ForwardThreshold = 10
)

func NewDomainInformationForwarder(forwardSignalChannel chan *models.ForwardMessage, routingInformationChannel chan *models.DomainController) *DomainInformationForwarder {
	forwarder := new(DomainInformationForwarder)
	forwarder.forwardSignalChannel = forwardSignalChannel
	forwarder.routingInformationChannel = routingInformationChannel
	forwarder.routingManager = forwardRouting.NewForwardRoutingManager(configuration.DomainControllerConfigurationManagerInstance().Config().ClusterManagementAddress)
	forwarder.forwardPriorityCounter = make(map[string]int)

	informationStorageDelegate,error := persistence.NewDomainInformationStorage()
	if error != nil {
		fmt.Println(error)
		os.Exit(1)
	}
	forwarder.informationStorageDelegate = informationStorageDelegate
	forwarder.forwarderStarted.Add(1)
	forwarder.forwarderStopped.Add(1)
	go forwarder.run()
	forwarder.forwarderStarted.Wait()
	return forwarder
}

//Starts listening on forward and routingInformation signals
func (forwarder *DomainInformationForwarder) run() {
	go forwarder.listenOnForwardSignal()
	go forwarder.listenOnRoutingInformationChannel()
	go forwarder.startForwardTicker()
	forwarder.forwarderStarted.Done()
}


//Listens on channel for ForwardMessage from DomainInformationProcessor
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

//Listens on channel for RoutingInformation from DomainInformationProcessor
func (forwarder *DomainInformationForwarder) listenOnRoutingInformationChannel() {
	for {
		newDomainController, open := <-forwarder.routingInformationChannel
		if !open {
			break
		}
		if newDomainController != nil {
			go forwarder.processNewDomainController(newDomainController  )
		}
	}
}


//Domain Controller information from DomainProcessor are added to Routing Manager cache
func (forwarder *DomainInformationForwarder) processNewDomainController(domainController *models.DomainController) {
	forwarder.routingManager.AddDomainControllerForDomain(domainController,domainController.Domain)
}

//Starts ticker which checks repeadetly if DomainInformation of a Real World Domain should be forwarded
func (forwarder *DomainInformationForwarder) startForwardTicker() {
	forwardTicker := time.NewTicker(DomainForwardInterval)
	for _ = range forwardTicker.C {
		forwarder.forwardsAllDomainInformation()
	}
}


//Proccesses forward messages from DomainInformationProcessor by forwarding all DomainInformationMessages for the corresponding Domain
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

//Forward DomainMessages for each Real World Domain
func (forwarder *DomainInformationForwarder) forwardsAllDomainInformation() {
	domains, _ := forwarder.informationStorageDelegate.FindAllDomains()
	for _, domain := range domains {
		forwarder.forwardDomainInformation(domain)
	}
}


//Forward DomainInformationMessages for a Real World Domain.
//Fetches messages from database.
//Request Parent Domain Controller.
//Delegates forwarding to MqttPublisher.
func (forwarder *DomainInformationForwarder) forwardDomainInformation(domain *models.RealWorldDomain) {
	forwarder.forwardPriorityCounter[domain.Name] = 0

	domainInformation, err := forwarder.informationStorageDelegate.FindDomainInformationByDomainName(domain.Name)

	if err != nil {
		fmt.Println(err)
		return
	}

	if len(domainInformation) == 0 {
		forwarder.informationStorageDelegate.RemoveDomain(domain)
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