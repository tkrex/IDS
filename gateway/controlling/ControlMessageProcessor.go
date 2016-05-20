package controlling

import (
	"sync"
	"github.com/tkrex/IDS/common/models"
	"sync/atomic"
	"fmt"
	"encoding/json"
	"github.com/tkrex/IDS/common/publishing"
	"github.com/tkrex/IDS/gateway/persistence"
)

type ControlMessageProcessor struct {
	workerStarted                  sync.WaitGroup
	workerStopped                  sync.WaitGroup
	state                          int64
	incomingControlMessagesChannel chan *models.ControlMessage
	informationPublisher           *publishing.MqttPublisher
}

func NewControlMessageProcessor(incomingControlMessagesChanel chan *models.ControlMessage) *ControlMessageProcessor {
	worker := new(ControlMessageProcessor)
	worker.workerStarted.Add(1)
	worker.workerStopped.Add(1)
	worker.incomingControlMessagesChannel = incomingControlMessagesChanel
	go worker.run()
	worker.workerStarted.Wait()
	fmt.Println("ControlMessageForwarder started")
	return worker
}

func (worker *ControlMessageProcessor) run() {
	PublishConfig := models.NewMqttClientConfiguration("tcp://localhost:1883", "ControlMessage", "gateway")

	worker.informationPublisher = publishing.NewMqttPublisher(PublishConfig)
	worker.workerStarted.Done()

	for closed := atomic.LoadInt64(&worker.state) == 1; !closed; closed = atomic.LoadInt64(&worker.state) == 1 {
		open := worker.processIncomingControlMessage()
		if !open {
			worker.Close()
			break
		}
	}
	worker.workerStopped.Done()
}

func (worker *ControlMessageProcessor) processIncomingControlMessage() bool {



	controlMessage, open := <-worker.incomingControlMessagesChannel
	 dbWorker := persistence.NewGatewayDBWorker()
	if dbWorker == nil {
		fmt.Println("Can't connect to database")
		return false
	}
	 defer dbWorker.Close()

	if controlMessage == nil {
		return false
	}

	if controlMessage.MessageType == models.DomainControllerDelete {
		dbWorker.RemoveDomainControllers(controlMessage.DomainControllers)
		worker.forwardControlMessage(controlMessage)
	} else if controlMessage.MessageType == models.DomainControllerDelete {
		forwardedDomainControllers := [] *models.DomainController{}
		for _, domainController := range controlMessage.DomainControllers {
			updated, _ := dbWorker.UpdateControllerInformation(domainController)
			if !updated {
				fmt.Println("Domain Controller information already exist")
			} else {
				forwardedDomainControllers = append(forwardedDomainControllers, domainController)
			}
		}
		if len(forwardedDomainControllers) > 0 {
			forwardControlMessage := models.NewControlMessage(controlMessage.MessageType, forwardedDomainControllers)
			worker.forwardControlMessage(forwardControlMessage)
		}
	}

	return open
}

func (worker *ControlMessageProcessor) forwardControlMessage(controlMessage *models.ControlMessage) {
	json, err := json.Marshal(&controlMessage)
	if err != nil {
		fmt.Print(err)
	} else {
		go worker.informationPublisher.Publish(json)
	}
}

func (worker *ControlMessageProcessor) Close() {
	fmt.Println("Closing ServerMaintenanceWorker")
	atomic.StoreInt64(&worker.state, 1)
	worker.workerStopped.Wait()
	fmt.Println("ServerMaintenanceWorker Closed")
}


