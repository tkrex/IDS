package providing

import (
	"github.com/tkrex/IDS/common/models"
	"github.com/tkrex/IDS/domainController/persistence"
)

type BrokerRequestHandler struct {

}

func NewBrokerRequestHandler() *BrokerRequestHandler {
	return new(BrokerRequestHandler)
}


func (handler *BrokerRequestHandler) handleRequest(informationRequest *models.DomainInformationRequest) ([]*models.Broker,error) {
	dbDelegate,err := persistence.NewDomainControllerDatabaseWorker()
	if err != nil {
		return nil, err
	}
	defer dbDelegate.Close()
	brokers,err := dbDelegate.FindBrokersForInformationRequest(informationRequest)
	return brokers, err
}

