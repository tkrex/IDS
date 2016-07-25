package providing

import (
	"github.com/tkrex/IDS/common/models"
	"time"
)

type DomainInformationForBrokerRequestHandler struct {

}

func NewDomainInformationForBrokerRequestHandler() *DomainInformationForBrokerRequestHandler {
	return new(DomainInformationForBrokerRequestHandler)
}

func (handler *DomainInformationForBrokerRequestHandler) handleRequest(informationRequest *models.DomainInformationRequest, brokerId string) (*models.DomainInformationMessage, error) {
	domain := models.NewRealWorldDomain(informationRequest.Domain())

	//DEBUG CODE
	broker := models.NewBroker()
	broker.ID = "testID"
	broker.IP = "localhost"
	broker.RealWorldDomain = models.NewRealWorldDomain("education")
	broker.Geolocation = models.NewGeolocation("Germany","Bavaria","Garching",11.6309,48.2499)

	topics := []*models.Topic{}

	for i := 0; i < 5; i++ {
		topic := models.NewTopic("/home/kitchen","{\"temperature\":3}",time.Now())
		topic.UpdateBehavior.UpdateIntervalDeviation = 3.0
		topic.PayloadSimilarity = 80.5
		topic.UpdateBehavior.Reliability = "automatic"
		topics = append(topics, topic)
	}

	message := models.NewDomainInformationMessage(domain,broker,topics)
	return message,nil

	//dbDelegate, err := persistence.NewDomainControllerDatabaseWorker()
	//if err != nil {
	//	fmt.Println(err)
	//	return nil, err
	//}
	//defer dbDelegate.Close()
	//domainInformation, err := dbDelegate.FindDomainInformationForBroker(informationRequest, brokerId)
	//return domainInformation, err
}

