package processing

import (
	"sync/atomic"
	"sync"
	"fmt"
	"encoding/json"
	"github.com/tkrex/IDS/common/models"
	"github.com/tkrex/IDS/domainController/persistence"
	"github.com/tkrex/IDS/domainController/configuration"
	"github.com/tkrex/IDS/domainController/scaling"
)

//Processes incoming DomainInformationMessages from Daemons and Sub Level Domain Controllers
type DomainInformationProcessor struct {
	// Flag to indicate that the consumer state
	// 0 == Running
	// 1 == Closed
	state                      int64
	processorStarted           sync.WaitGroup
	processorStopped           sync.WaitGroup
	incomingInformationChannel chan *models.RawTopicMessage
	forwardFlag                bool
	forwardingSignalChannel    chan *models.ForwardMessage
	routingInformationChannel  chan *models.DomainController
	scalingManager             *scaling.ScalingRequestManager


}

func NewDomainInformationProcessor(incomingInformationChannel chan *models.RawTopicMessage, forwardFlag bool) *DomainInformationProcessor {
	processor := new(DomainInformationProcessor)
	processor.processorStarted.Add(1)
	processor.processorStopped.Add(1)
	processor.incomingInformationChannel = incomingInformationChannel
	processor.forwardFlag = forwardFlag
	processor.forwardingSignalChannel = make(chan *models.ForwardMessage)
	processor.routingInformationChannel = make(chan *models.DomainController)
	processor.scalingManager = scaling.NewScalingManager(configuration.DomainControllerConfigurationManagerInstance().Config().ClusterManagementAddress)
	go processor.run()
	processor.processorStarted.Wait()
	fmt.Println("Producer Created")
	return processor
}


//Returns channel via which the DomainInformationProcessor signals the forwarding of DomainInformationMessages
func (processor *DomainInformationProcessor) ForwardSignalChannel() chan *models.ForwardMessage {
	return processor.forwardingSignalChannel
}

//Returns channel which is used to send Domain Controller information to DomainInformationForwarder
func (processor *DomainInformationProcessor) RoutingInformationChannel() chan *models.DomainController {
	return processor.routingInformationChannel
}

//Stops DomainInformation Processors
func (processor *DomainInformationProcessor)  Close() {
	fmt.Println("Closing Processor")
	atomic.StoreInt64(&processor.state, 1)
	processor.processorStopped.Wait()
	fmt.Println("Processor Closed")
}


//Starts processing incoming DomainInformationMessages
func (processor *DomainInformationProcessor) run() {

	processor.processorStarted.Done()
	for closed := atomic.LoadInt64(&processor.state) == 1; !closed; closed = atomic.LoadInt64(&processor.state) == 1 {
		open := processor.processDomainInformationMessages()
		if !open {
			processor.Close()
			break
		}
	}
	processor.processorStopped.Done()
}

//Processes incoming DomainInformationMessages received from a MqttSubscriber
func (processor *DomainInformationProcessor) processDomainInformationMessages() bool {
	rawTopic, ok := <-processor.incomingInformationChannel
	if rawTopic != nil {
		go processor.processDomainInformationMessage(rawTopic)
	}
	return ok
}

///Processes a RawTopicMessag containing a DomainInformationMessage
//Checks if DomainInformationMessages for a Real World Domain should be forwarded
func (processor *DomainInformationProcessor) processDomainInformationMessage(topic *models.RawTopicMessage) {
	var domainInformationMessage *models.DomainInformationMessage
	if err := json.Unmarshal(topic.Payload, &domainInformationMessage); err != nil {
		fmt.Println(err)
		return
	}

	processor.storeDomainInformation(domainInformationMessage)

	if processor.forwardFlag {
		processor.forwardingSignalChannel <- models.NewForwardMessage(domainInformationMessage.RealWorldDomain,domainInformationMessage.ForwardPriority)
	}
	if processor.scalingManager.CheckWorkloadForDomain(domainInformationMessage.RealWorldDomain) {
		if domainController := processor.scalingManager.CreateNewDominControllerForDomain(domainInformationMessage.RealWorldDomain); domainController != nil {
			processor.routingInformationChannel <- domainController
		}

	}
}

//Stores received DomainInformationMessages in database
func (processor *DomainInformationProcessor) storeDomainInformation(information *models.DomainInformationMessage) {
	dbDelegate, _ := persistence.NewDomainInformationStorage()
	if dbDelegate == nil {
		return
	}
	defer dbDelegate.Close()
	dbDelegate.StoreDomainInformationMessage(information)
}
