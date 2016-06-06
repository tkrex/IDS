package controlling

import (
	"sync"
	"github.com/tkrex/IDS/common/models"
	"fmt"
	"sync/atomic"
	"encoding/json"
)

type ControlMessageProcessor struct {
	// Flag to indicate that the consumer state
	// 0 == Running
	// 1 == Closed
	state                         int64
	processorStarted              sync.WaitGroup
	processorStopped              sync.WaitGroup
	incomingControlMessageChannel chan *models.RawTopicMessage
}


func NewControlMessageProcessor(incomingControlMessageChannel chan *models.RawTopicMessage) *ControlMessageProcessor {
	processor := new(ControlMessageProcessor)
	processor.processorStarted.Add(1)
	processor.processorStopped.Add(1)
	processor.incomingControlMessageChannel = incomingControlMessageChannel

	go processor.run()
	processor.processorStarted.Wait()
	fmt.Println("Processor Created")
	return processor
}


func (processor *ControlMessageProcessor) run() {
	processor.processorStarted.Done()

	for closed := atomic.LoadInt64(&processor.state) == 1; !closed; closed = atomic.LoadInt64(&processor.state) == 1 {
		open := processor.processIncomingControlMessage()
		if !open {
			processor.Close()
			break
		}
	}
	processor.processorStopped.Done()
}

func (processor *ControlMessageProcessor) processIncomingControlMessage() bool {
	message, open := <- processor.incomingControlMessageChannel
	fmt.Println("Received ControlMessage")
	if message.Name != "ControlMessage" {
		fmt.Println("Wrong topic name")
		return open
	}
	var controlMessage *models.ControlMessage
	err := json.Unmarshal(message.Payload,&controlMessage)
	if err != nil {
		fmt.Println(err)
		return open
	}

	if controlMessage != nil {
		if controlMessage.MessageType == models.DomainControllerChange {
			go processor.storeDomainController(controlMessage.DomainController)
		} else if controlMessage.MessageType == models.DomainControllerDelete {
			go processor.removeDomainController(controlMessage.DomainController)
		}
	}
	return open
}


func (processor *ControlMessageProcessor) storeDomainController(controller *models.DomainController) {
	dbWorker, err := NewControlMessageDBDelegate()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dbWorker.Close()
	dbWorker.StoreDomainController(controller)
}

func (processor *ControlMessageProcessor) removeDomainController(controller *models.DomainController) {
	dbWorker, err := NewControlMessageDBDelegate()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dbWorker.Close()
	dbWorker.RemoveDomainControllerForDomain(controller.Domain)
}

func (processor *ControlMessageProcessor)  Close() {
	fmt.Println("Closing ControlMessageProcessor")
	atomic.StoreInt64(&processor.state, 1)
	processor.processorStopped.Wait()
	fmt.Println("Processor ControlMessageProcessor")
}
