package providing

import (
	"net/http"
	"github.com/tkrex/IDS/common/controlling"
	"github.com/tkrex/IDS/common/models"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"time"
)

type DomainInformationForBrokerRequestHandler struct {

}

func NewDomainInformationForBrokerRequestHandler() *DomainInformationForBrokerRequestHandler {
	return new(DomainInformationForBrokerRequestHandler)
}


func (handler *DomainInformationForBrokerRequestHandler) handleRequest(brokerId string) *models.DomainInformationMessage {

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
	return message

	dbDelegate, err := controlling.NewControlMessageDBDelegate()
	if err != nil {
		return nil
	}
	defer dbDelegate.Close()

	var destinationDomainController *models.DomainController
	domain = models.NewRealWorldDomain(brokerId)
	destinationDomainController = dbDelegate.FindDomainControllerForDomain(domain.FirstLevelDomain())

	if destinationDomainController == nil {
		destinationDomainController = dbDelegate.FindDomainControllerForDomain("default")
	}

	if destinationDomainController != nil {
		return handler.forwardRequestToDomainController(brokerId,destinationDomainController)
	}
	return nil
}


//TODO: Add Endpoint at DomainController
func (handler *DomainInformationForBrokerRequestHandler) forwardRequestToDomainController(domainName string, domainController *models.DomainController) *models.DomainInformationMessage {
	requestUrl := domainController.RestEndpoint.String() + "/domainController/domainInformation/" + domainName
	fmt.Println("Forwarding Request to ",requestUrl)
	client := http.DefaultClient
	response, err := client.Get(requestUrl)
	if err != nil {
		fmt.Printf("%s", err)
		return nil
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("%s", err)
		return nil
	}
	var domainInformation *models.DomainInformationMessage
	json.Unmarshal([]byte(string(contents)), &domainInformation)
	return domainInformation
}