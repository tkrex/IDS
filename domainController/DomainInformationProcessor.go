package domainController

import (
	"sync/atomic"
	"sync"
	"fmt"
	"encoding/json"
	"github.com/tkrex/IDS/common/models"
)

type DomainInformationProcessor struct {
	// Flag to indicate that the consumer state
	// 0 == Running
	// 1 == Closed
	state                               int64
	processorStarted                    sync.WaitGroup
	processorStopped                    sync.WaitGroup
	incomingTopicChannel                chan *models.RawTopicMessage
	domainInformationPersistenceManager DomainInformationPersistenceManager
}

func NewDomainInformationProcessor(persistenceManager DomainInformationPersistenceManager, incomingTopicChannel chan *models.RawTopicMessage) *DomainInformationProcessor {
	processor := new(DomainInformationProcessor)
	processor.processorStarted.Add(1)
	processor.processorStopped.Add(1)
	processor.domainInformationPersistenceManager = persistenceManager
	processor.incomingTopicChannel = incomingTopicChannel
	go processor.run()
	processor.processorStarted.Wait()
	fmt.Println("Producer Created")
	return processor
}

func (processor *DomainInformationProcessor) State() int64 {
	return atomic.LoadInt64(&processor.state)
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
		open := processor.ProcessDomainInformationMessages()
		if !open {
			processor.Close()
			break
		}
	}
	processor.processorStopped.Done()
}

func (processor *DomainInformationProcessor) ProcessDomainInformationMessages() bool {
	rawTopic, ok := <-processor.incomingTopicChannel
	if rawTopic != nil {
		processor.processDomainInformationMessage(rawTopic)
	}
	return ok
}

func (processor *DomainInformationProcessor) processDomainInformationMessage(topic *models.RawTopicMessage) {
	var domainInformationMessage models.DomainInformationMessage
	if err := json.Unmarshal(topic.Payload, &domainInformationMessage); err != nil {
		fmt.Println(err)
		return
	}
	processor.domainInformationPersistenceManager.Store(&domainInformationMessage)
}
