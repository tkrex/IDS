package layers

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
		open := processor.processIncomingControlMessages()
		if !open {
			processor.Close()
			break
		}
	}
	processor.processorStopped.Done()
}

func (processor *ControlMessageProcessor) processIncomingControlMessages() bool {
	message, open := <- processor.incomingControlMessageChannel
	fmt.Println("Received ControlMessage")
	if message.Name != "ControlMessage" {
		fmt.Println("Wrong topic name")
		return open
	}
	var domainControllers []*models.DomainController
	err := json.Unmarshal(message.Payload,&domainControllers)
	if err != nil {
		fmt.Println(err)
		return open
	}

	if domainControllers != nil {
		go processor.storeDomainControllers(domainControllers)
	}
	return open
}


func (processor *ControlMessageProcessor) storeDomainControllers(controllers []*models.DomainController) {
	dbWorker, err := NewDaemonDatabaseWorker()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dbWorker.Close()
	dbWorker.StoreDomainControllers(controllers)
}

func (processor *ControlMessageProcessor)  Close() {
	fmt.Println("Closing ControlMessageProcessor")
	atomic.StoreInt64(&processor.state, 1)
	processor.processorStopped.Wait()
	fmt.Println("Processor ControlMessageProcessor")
}
