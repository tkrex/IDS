package controlling

import (
	"github.com/tkrex/IDS/common/models"
	"fmt"
	"encoding/json"
	"github.com/tkrex/IDS/common/publishing"
	"github.com/tkrex/IDS/common/controlling"
	"net/url"
)

type DomainControllerManagementRequestHandler struct {
	managementBrokerAddress string
}

func NewDomainControllerManagementRequestHandler(managementBrokerAddress string) *DomainControllerManagementRequestHandler {
	worker := new(DomainControllerManagementRequestHandler)
	worker.managementBrokerAddress = managementBrokerAddress
	return worker
}

func (handler *DomainControllerManagementRequestHandler) handleManagementRequest(request *models.DomainControllerManagementRequest) *models.DomainController {
	 dbWorker,_ := controlling.NewControlMessageDBDelegate()
	if dbWorker == nil {
		fmt.Println("Can't connect to database")
		return nil
	}
	 defer dbWorker.Close()

	var changedDomainController  *models.DomainController
	if request.RequestType == models.DomainControllerDelete {
		if domainController := dbWorker.FindDomainControllerForDomain(request.Domain.Name); domainController != nil {
			dbWorker.RemoveDomainControllerForDomain(request.Domain)
			changedDomainController = domainController
		}

	} else if request.RequestType == models.DomainControllerChange {
		if domainController := handler.startNewDomainControllerInstance(request.Domain); domainController != nil {
			changed, _ := dbWorker.StoreDomainController(domainController)
			if changed {
				fmt.Println("New Domain Controller stored")
				changedDomainController = domainController
			}
		}
	}

	if changedDomainController != nil {
		controlMessage := models.NewControlMessage(request.RequestType,changedDomainController)
		handler.forwardControlMessage(controlMessage)
	}
	return changedDomainController
}

func (handler *DomainControllerManagementRequestHandler) startNewDomainControllerInstance(domain *models.RealWorldDomain) *models.DomainController {
	//TODO: start new docker instance

	restEndpoint,_ := url.Parse("http://localhost:8000/rest")
	brokerAddress,_ := url.Parse("ws://localhost:11883")
	domainController := models.NewDomainController(restEndpoint,brokerAddress,domain)
	return domainController
}



func (worker *DomainControllerManagementRequestHandler) forwardControlMessage(controlMessage *models.ControlMessage) {
	json, err := json.Marshal(&controlMessage)
	if err != nil {
		fmt.Print(err)
		return
	}

	publishConfig := models.NewMqttClientConfiguration(worker.managementBrokerAddress,"1883","tcp", "ControlMessage", "gateway")
	publisher := publishing.NewMqttPublisher(publishConfig,false)
	defer publisher.Close()
	publisher.Publish(json)
}



