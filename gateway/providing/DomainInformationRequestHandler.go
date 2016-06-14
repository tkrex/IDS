package providing

import (
	"net/http"
	"github.com/tkrex/IDS/common/models"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"time"
	"github.com/tkrex/IDS/common/routing"
)

type DomainInformationRequestHandler struct {

}

func NewDomainInformationRequestHandler() *DomainInformationRequestHandler {
	return new(DomainInformationRequestHandler)
}


func (handler *DomainInformationRequestHandler) handleRequest(domainName string) []*models.DomainInformationMessage {
	domainInformation := []*models.DomainInformationMessage{}

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
	domainInformation = append(domainInformation,message)
	return domainInformation


	domain = models.NewRealWorldDomain(domainName)
	destinationController := routing.NewRoutingManager().DomainControllerForDomain(domain)

	if destinationController != nil {
		return handler.requestDomainInformationFromDomainController(domainName,destinationController)
	}
	return nil
}


func (handler *DomainInformationRequestHandler) requestDomainInformationFromDomainController(domainName string, domainController *models.DomainController) []*models.DomainInformationMessage {
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
	domainInformation := []*models.DomainInformationMessage{}
	json.Unmarshal([]byte(string(contents)), &domainInformation)
	return domainInformation
}