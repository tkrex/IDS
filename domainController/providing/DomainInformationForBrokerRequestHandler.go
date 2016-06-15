package providing

import (
	"github.com/tkrex/IDS/common/models"
	"fmt"
	"github.com/tkrex/IDS/domainController/persistence"
)

type DomainInformationForBrokerRequestHandler struct {

}

func NewDomainInformationForBrokerRequestHandler() *DomainInformationForBrokerRequestHandler {
	return new(DomainInformationForBrokerRequestHandler)
}

func (handler *DomainInformationForBrokerRequestHandler) handleRequest(informationRequest *models.DomainInformationRequest, brokerId string) ([]*models.DomainInformationMessage, error) {

	dbDelegate, err := persistence.NewDomainControllerDatabaseWorker()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer dbDelegate.Close()
	domainInformation, err := dbDelegate.FindDomainInformationForBroker(informationRequest, brokerId)
	return domainInformation, err
}

