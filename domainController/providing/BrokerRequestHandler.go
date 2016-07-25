package providing

import (
	"github.com/tkrex/IDS/common/models"
)

type BrokerRequestHandler struct {

}

func NewBrokerRequestHandler() *BrokerRequestHandler {
	return new(BrokerRequestHandler)
}


func (handler *BrokerRequestHandler) handleRequest(informationRequest *models.DomainInformationRequest) ([]*models.Broker,error) {

	broker1 := models.NewBroker()
	broker1.ID = "weatherBroker"
	broker1.IP = "12.12.12.12:1833"
	broker1.InternetDomain = "krex.in.tum.de"
	broker1.Statitics.NumberOfTopics = 1022
	broker1.Statitics.ReceivedTopicsPerSeconds = 10
	broker1.RealWorldDomain = models.NewRealWorldDomain("weather")
	broker1.Geolocation = models.NewGeolocation("Germany", "Bavaria", "Garching", 11.6309, 48.2499)

	broker := models.NewBroker()
	broker.ID = "weatherBroker"
	broker.IP = "143.142.192.192:1833"
	broker.InternetDomain = "unknown.de"
	broker.Statitics.NumberOfTopics = 1022
	broker.Statitics.ReceivedTopicsPerSeconds = 10
	broker.RealWorldDomain = models.NewRealWorldDomain("weather")
	broker.Geolocation = models.NewGeolocation("Germany", "NRW", "Siegen", 11.6309, 48.2499)
	return []*models.Broker{broker,broker1}, nil

	//dbDelegate,err := persistence.NewDomainControllerDatabaseWorker()
	//if err != nil {
	//	return nil, err
	//}
	//defer dbDelegate.Close()
	//brokers,err = dbDelegate.FindBrokersForInformationRequest(informationRequest)
	//return brokers, err
}

