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

	domains := []*models.RealWorldDomain{ models.NewRealWorldDomain("education/schools"), models.NewRealWorldDomain("education"), models.NewRealWorldDomain("education/university")}
	brokers := []*models.Broker{}
	for _, domain := range domains {
		//DEBUG CODE
		broker := models.NewBroker()
		broker.ID = "testID"
		broker.IP = "12.123.123.12:1883"
		broker.InternetDomain="www.krexit.co"
		broker.RealWorldDomain = domain
		broker.Geolocation = models.NewGeolocation("germany", "bavaria", "munich", 11.6309, 48.2499)
		broker.Statitics.NumberOfTopics = 20
		brokers = append(brokers, broker)
	}
	return brokers, nil

	dbDelegate,err := persistence.NewDomainControllerDatabaseWorker()
	if err != nil {
		return nil, err
	}
	defer dbDelegate.Close()
	brokers,err = dbDelegate.FindBrokersForInformationRequest(informationRequest)
	return brokers, err
}

