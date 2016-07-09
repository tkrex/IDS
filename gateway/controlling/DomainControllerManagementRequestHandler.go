package controlling

import (
	"github.com/tkrex/IDS/common/models"
	"fmt"
	"encoding/json"
	"github.com/tkrex/IDS/common/publishing"
	"github.com/tkrex/IDS/common/controlling"
	"net/url"
	"github.com/tkrex/IDS/gateway/scaling"
)

type DomainControllerManagementRequestHandler struct {
	managementBrokerAddress *url.URL
}

func NewDomainControllerManagementRequestHandler(managementBrokerAddress *url.URL) *DomainControllerManagementRequestHandler {
	worker := new(DomainControllerManagementRequestHandler)
	worker.managementBrokerAddress = managementBrokerAddress
	return worker
}

func (handler *DomainControllerManagementRequestHandler) handleManagementRequest(request *models.DomainControllerManagementRequest) *models.DomainController {
	dbWorker, _ := controlling.NewControlMessageDBDelegate()
	if dbWorker == nil {
		fmt.Println("Can't connect to database")
		return nil
	}
	defer dbWorker.Close()

	var changedDomainController  *models.DomainController
	switch request.RequestType {
	case models.DomainControllerDelete:
		if domainController := dbWorker.FindDomainControllerForDomain(request.Domain); domainController != nil {
			dbWorker.RemoveDomainControllerForDomain(request.Domain)
			changedDomainController = domainController
		}

	case models.DomainControllerChange:
		if domainController,_ := handler.startNewDomainControllerInstance(request.Domain); domainController != nil {
			changed, _ := dbWorker.StoreDomainController(domainController)
			if changed {
				fmt.Println("New Domain Controller stored")
				changedDomainController = domainController
			}
		}
	case models.DomainControllerFetch:
		dbDelegate , err := controlling.NewControlMessageDBDelegate()
		if err != nil {
			fmt.Println(err)
		}
		changedDomainController = dbDelegate.FindDomainControllerForDomain(request.Domain)
	}


	if changedDomainController != nil {
		controlMessage := models.NewControlMessage(request.RequestType, changedDomainController)
		handler.forwardControlMessage(controlMessage)
	}
	return changedDomainController
}

func (handler *DomainControllerManagementRequestHandler) startNewDomainControllerInstance(domain *models.RealWorldDomain) (*models.DomainController,error) {
	//TODO: start new docker instance
	domainController,error := scaling.NewScalingManager().StartDomainControllerInstance(domain)
	return domainController, error
}

func (worker *DomainControllerManagementRequestHandler) forwardControlMessage(controlMessage *models.ControlMessage) {
	json, err := json.Marshal(&controlMessage)
	if err != nil {
		fmt.Print(err)
		return
	}

	publishConfig := models.NewMqttClientConfiguration(worker.managementBrokerAddress, "gateway")
	publisher := publishing.NewMqttPublisher(publishConfig, false)
	defer publisher.Close()
	publisher.Publish(json, "ControlMessage")
}



