package processing

import (
	"sync/atomic"
	"sync"
	"fmt"
	"encoding/json"
	"github.com/tkrex/IDS/common/models"
	"github.com/tkrex/IDS/domainController/persistence"
	"github.com/tkrex/IDS/domainController/configuration"
	"os"
	"github.com/tkrex/IDS/common/routing"
	"github.com/tkrex/IDS/domainController"
)







type DomainInformationProcessor struct {
	// Flag to indicate that the consumer state
	// 0 == Running
	// 1 == Closed
	state                       int64
	processorStarted            sync.WaitGroup
	processorStopped            sync.WaitGroup
	incomingTopicChannel        chan *models.RawTopicMessage
	forwardFlag bool
	forwardingSignalChannel     chan *models.ForwardMessage
	routingManager 	*routing.RoutingManager
	scalingManager  *domainController.ScalingManager


}

func NewDomainInformationProcessor(incomingTopicChannel chan *models.RawTopicMessage, forwardFlag bool) *DomainInformationProcessor {
	processor := new(DomainInformationProcessor)
	processor.processorStarted.Add(1)
	processor.processorStopped.Add(1)
	processor.incomingTopicChannel = incomingTopicChannel
	processor.forwardFlag = forwardFlag
	processor.forwardingSignalChannel = make(chan *models.RealWorldDomain)

	processor.routingManager = routing.NewRoutingManager(configuration.DomainControllerConfigurationManagerInstance().Config().ScalingInterfaceAddress)
	processor.scalingManager = domainController.NewScalingManager()
	go processor.run()
	processor.processorStarted.Wait()
	fmt.Println("Producer Created")
	return processor
}

func (processor *DomainInformationProcessor) State() int64 {
	return atomic.LoadInt64(&processor.state)
}

func (processor *DomainInformationProcessor) ForwardSignalChannel() chan *models.ForwardMessage {
	return processor.forwardingSignalChannel
}

func (processor *DomainInformationProcessor)  Close() {
	fmt.Println("Closing Processor")
	atomic.StoreInt64(&processor.state, 1)
	processor.processorStopped.Wait()
	fmt.Println("Processor Closed")
}

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

func (processor *DomainInformationProcessor) processDomainInformationMessages() bool {
	rawTopic, ok := <-processor.incomingTopicChannel
	if rawTopic != nil {
		go processor.processDomainInformationMessage(rawTopic)
	}
	return ok
}


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
			processor.routingManager.AddDomainControllerForDomain(domainController, domainInformationMessage.RealWorldDomain)
		}

	}
}

func (processor *DomainInformationProcessor) storeDomainInformation(information *models.DomainInformationMessage) {
	dbDelegate, _ := persistence.NewDomainControllerDatabaseWorker()
	if dbDelegate == nil {
		return
	}
	defer dbDelegate.Close()
	dbDelegate.StoreDomainInformationMessage(information)
}
