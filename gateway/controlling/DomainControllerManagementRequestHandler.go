package controlling

import (
	"github.com/tkrex/IDS/common/models"
	"fmt"
	"encoding/json"
	"github.com/tkrex/IDS/common/publishing"
	"github.com/tkrex/IDS/gateway/persistence"
)

type DomainControllerManagementRequestHandler struct {

}

func NewDomainControllerManagementRequestHandler() *DomainControllerManagementRequestHandler {
	worker := new(DomainControllerManagementRequestHandler)
	return worker
}

func (handler *DomainControllerManagementRequestHandler) handleManagementRequest(request *models.DomainControllerManagementRequest) *models.DomainController {
	 dbWorker := persistence.NewGatewayDBWorker()
	if dbWorker == nil {
		fmt.Println("Can't connect to database")
		return nil
	}
	 defer dbWorker.Close()

	var changedDomainController  *models.DomainController
	if request.RequestType == models.DomainControllerDelete {
		if domainController := dbWorker.FindDomainControllerForDomain(request.Domain); domainController != nil {
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
	domainController := models.NewDomainController("localhost",domain)
	return domainController
}



func (worker *DomainControllerManagementRequestHandler) forwardControlMessage(controlMessage *models.ControlMessage) {
	json, err := json.Marshal(&controlMessage)
	if err != nil {
		fmt.Print(err)
		return
	}

	publishConfig := models.NewMqttClientConfiguration("localhost", "ControlMessage", "gateway")
	publisher := publishing.NewMqttPublisher(publishConfig)
	defer publisher.Close()
	publisher.Publish(json)

}


