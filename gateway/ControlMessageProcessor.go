package gateway

import (
	"sync"
	"github.com/tkrex/IDS/common/layers"
	"github.com/tkrex/IDS/common/models"
	"sync/atomic"
	"encoding/json"
	"fmt"
)


type ControlMessageProcessor struct {
	workerStarted sync.WaitGroup
	workerStopped sync.WaitGroup
	state int64
	webInterface *ServerMaintenanceWebInterface
	informationPublisher *common.MqttPublisher
}

func NewControlMessageForwarder() *ControlMessageProcessor {
	worker := new(ControlMessageProcessor)
	worker.workerStarted.Add(1)
	worker.workerStopped.Add(1)
	go worker.run()
	worker.workerStarted.Wait()
	fmt.Println("ControlMessageForwarder started")
	return worker
}

func (worker *ControlMessageProcessor) run() {
	PublishConfig := models.NewMqttClientConfiguration("tcp://localhost:1883","controlInformation","gateway")

	worker.informationPublisher = common.NewMqttPublisher(PublishConfig)
	worker.webInterface = NewServerMaintenanceWebInterface("8080")
	worker.workerStarted.Done()

	for closed := atomic.LoadInt64(&worker.state) == 1; !closed; closed = atomic.LoadInt64(&worker.state) == 1 {
		open := worker.processIncomingDomainControllerInformation()
		if !open {
			worker.Close()
			break
		}
	}
	worker.workerStopped.Done()
}


func (worker *ControlMessageProcessor) processIncomingDomainControllerInformation() bool {
	domainControllerInformation , open := <-worker.webInterface.incomingControlMessagesChannel
	if domainControllerInformation != nil {
		updated , _ := UpdateControllerInformation(domainControllerInformation)
		if !updated {
			fmt.Println("Domain COntroller information already exist")
			return open
		}
		json, err := json.Marshal(&domainControllerInformation)
		if err != nil {
			fmt.Print(err)
		} else {
			go worker.informationPublisher.Publish(json)
		}
	}
	return open
}

func (worker *ControlMessageProcessor) Close() {
	fmt.Println("Closing ServerMaintenanceWorker")
	atomic.StoreInt64(&worker.state, 1)
	worker.workerStopped.Wait()
	fmt.Println("ServerMaintenanceWorker Closed")
}


