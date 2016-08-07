package providing

import (
	"github.com/tkrex/IDS/common/models"
	"fmt"
	"time"
	"github.com/tkrex/IDS/domainController/persistence"
)

type DomainInformationRequestHandler struct {

}

func NewDomainInformationRequestHandler() *DomainInformationRequestHandler {
	return new(DomainInformationRequestHandler)
}


func (handler *DomainInformationRequestHandler) HandleRequest(informationRequest *models.DomainInformationRequest) ([]*models.DomainInformationMessage,error) {
	demo := []*models.DomainInformationMessage{}

	//DEBUG CODE
	domain := models.NewRealWorldDomain("education")
	broker := models.NewBroker()
	broker.ID = "testID"
	broker.IP = "localhost"
	broker.RealWorldDomain = models.NewRealWorldDomain("education")
	broker.Geolocation = models.NewGeolocation("germany","bavaria","munich",11.6309,48.2499)

	topics := []*models.TopicInformation{}

	for i := 0; i < 5; i++ {
		topic := models.NewTopicInformation("/home/kitchen","{\"temperature\":3}",time.Now())
		topic.UpdateBehavior.UpdateIntervalDeviation = 3.0
		topics = append(topics, topic)
	}

	message := models.NewDomainInformationMessage(domain,broker,topics)
	demo = append(demo,message)
	return demo, nil

	dbDelegate,err := persistence.NewDomainInformationStorage()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer dbDelegate.Close()
	domainInformation, err := dbDelegate.FindDomainInformationForRequest(informationRequest)
	return domainInformation, err
}

