package providing

import (
	"net/http"
	"github.com/tkrex/IDS/common/models"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"time"
	"github.com/tkrex/IDS/common/routing"
	"github.com/tkrex/IDS/domainController/persistence"
)

type DomainInformationRequestHandler struct {

}

func NewDomainInformationRequestHandler() *DomainInformationRequestHandler {
	return new(DomainInformationRequestHandler)
}


func (handler *DomainInformationRequestHandler) handleRequest(informationRequest *models.DomainInformationRequest) ([]*models.DomainInformationMessage,error) {
	demo := []*models.DomainInformationMessage{}

	//DEBUG CODE
	domain := models.NewRealWorldDomain("education")
	broker := models.NewBroker()
	broker.ID = "testID"
	broker.IP = "localhost"
	broker.RealWorldDomain = models.NewRealWorldDomain("education")
	broker.Geolocation = models.NewGeolocation("germany","bavaria","munich",11.6309,48.2499)

	topics := []*models.Topic{}

	for i := 0; i < 5; i++ {
		topic := models.NewTopic("/home/kitchen","{\"temperature\":3}",time.Now())
		topic.UpdateBehavior.UpdateIntervalDeviation = 3.0
		topics = append(topics, topic)
	}

	message := models.NewDomainInformationMessage(domain,broker,topics)
	demo = append(demo,message)
	return demo

	dbDelegate,err := persistence.NewDomainControllerDatabaseWorker()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer dbDelegate.Close()
	domainInformation, err := dbDelegate.FindDomainInformationForRequest(informationRequest)
	return domainInformation, err
}

