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

	broker := models.NewBroker()
	topics := []*models.Topic{}

	broker = models.NewBroker()
	broker.ID = "weatherBroker"
	broker.IP = "12.12.12.12:1833"
	broker.InternetDomain = "krex.in.tum.de"
	broker.Statitics.NumberOfTopics = 1022
	broker.Statitics.ReceivedTopicsPerSeconds = 10
	broker.RealWorldDomain = models.NewRealWorldDomain("weather")
	broker.Geolocation = models.NewGeolocation("Germany", "Bavaria", "Garching", 11.6309, 48.2499)

	topic := models.NewTopic("/fmi/server-room", "{\"temperature\":-6}", time.Now())
	topic.UpdateBehavior.AverageUpdateIntervalInSeconds = 180
	topic.UpdateBehavior.UpdateIntervalDeviation = 3.0
	topic.PayloadSimilarity = 80.5
	topic.UpdateBehavior.Reliability = "automatic"
	topics = append(topics, topic)

	topic = models.NewTopic("/fmi/ls1", "{\"temperature\":30}", time.Now())
	topic.UpdateBehavior.AverageUpdateIntervalInSeconds = 120
	topic.UpdateBehavior.UpdateIntervalDeviation = 70.0
	topic.PayloadSimilarity = 87
	topic.UpdateBehavior.Reliability = "semi-automatic"
	topics = append(topics, topic)

	topic = models.NewTopic("/fmi/ls2", "{\"temperature\":10}", time.Now())
	topic.UpdateBehavior.AverageUpdateIntervalInSeconds = 30
	topic.UpdateBehavior.UpdateIntervalDeviation = 200
	topic.PayloadSimilarity = 40
	topic.UpdateBehavior.Reliability = "non-deterministic"
	topics = append(topics, topic)



	message := models.NewDomainInformationMessage(domain, broker, topics)
	return message, nil

	//dbDelegate, err := persistence.NewDomainControllerDatabaseWorker()
	//if err != nil {
	//	fmt.Println(err)
	//	return nil, err
	//}
	//defer dbDelegate.Close()
	//domainInformation, err := dbDelegate.FindDomainInformationForBroker(informationRequest, brokerId)
	//return domainInformation, err
}

