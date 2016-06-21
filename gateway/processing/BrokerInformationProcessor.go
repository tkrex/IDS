package processing

import (
	"sync/atomic"
	"sync"
	"fmt"
	"encoding/json"
	"github.com/tkrex/IDS/common/models"
	"github.com/tkrex/IDS/gateway/persistence"
)

const (
	ForwardThreshold = 10
)

type DomainProcessor struct {
	// Flag to indicate that the consumer state
	// 0 == Running
	// 1 == Closed
	state                  int64
	processorStarted       sync.WaitGroup
	processorStopped       sync.WaitGroup
	incomingBrokersChannel chan *models.RawTopicMessage

}

func NewDomainProcessor(incomingTopicChannel chan *models.RawTopicMessage) *DomainProcessor {
	processor := new(DomainProcessor)
	processor.processorStarted.Add(1)
	processor.processorStopped.Add(1)
	processor.incomingBrokersChannel = incomingTopicChannel
	go processor.run()
	processor.processorStarted.Wait()
	fmt.Println("Producer Created")
	return processor
}

func (processor *DomainProcessor) State() int64 {
	return atomic.LoadInt64(&processor.state)
}

func (processor *DomainProcessor)  Close() {
	fmt.Println("Closing Processor")
	atomic.StoreInt64(&processor.state, 1)
	processor.processorStopped.Wait()
	fmt.Println("Processor Closed")
}

func (processor *DomainProcessor) run() {

	processor.processorStarted.Done()
	for closed := atomic.LoadInt64(&processor.state) == 1; !closed; closed = atomic.LoadInt64(&processor.state) == 1 {
		open := processor.listenForDomains()
		if !open {
			processor.Close()
			break
		}
	}
	processor.processorStopped.Done()
}

func (processor *DomainProcessor) listenForDomains() bool {
	rawTopic, ok := <-processor.incomingBrokersChannel
	if rawTopic != nil {
		go processor.processDomainList(rawTopic)
	}
	return ok
}

func (processor *DomainProcessor) processDomainList(topic *models.RawTopicMessage) {
	var brokers []*models.Broker
	if err := json.Unmarshal(topic.Payload, &brokers); err != nil {
		fmt.Println(err)
		return
	}

	processor.storeDomainList(brokers)
}

func (processor *DomainProcessor) storeDomainList(domains []*models.RealWorldDomain) {
	dbDelegate := persistence.NewGatewayDBWorker()
	if dbDelegate == nil {
		return
	}
	defer dbDelegate.Close()
	dbDelegate.S
}