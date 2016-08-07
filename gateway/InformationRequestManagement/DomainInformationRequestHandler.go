package informationRequestManagement

import (
	"net/http"
	"github.com/tkrex/IDS/common/models"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"time"
	"net/url"
)

type DomainInformationRequestHandler struct {

}

func NewDomainInformationRequestHandler() *DomainInformationRequestHandler {
	return new(DomainInformationRequestHandler)
}


func (handler *DomainInformationRequestHandler) HandleRequest(informationRequest *models.DomainInformationRequest) []*models.DomainInformationMessage {
	domainInformation := []*models.DomainInformationMessage{}

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
	domainInformation = append(domainInformation,message)
	return domainInformation


	domain = models.NewRealWorldDomain(informationRequest.Domain())
	destinationController := RequestRoutingManagerInstance().DomainControllerForDomain(domain)

	if destinationController != nil {
		return handler.requestDomainInformationFromDomainController(informationRequest,destinationController)
	}
	return nil
}


func (handler *DomainInformationRequestHandler) requestDomainInformationFromDomainController(informationRequest *models.DomainInformationRequest, domainController *models.DomainController) []*models.DomainInformationMessage {
	requestUrlString := domainController.RestEndpoint.String() + "/domainController/domainInformation/" + informationRequest.Domain()
	requestUrl,_ := url.Parse(requestUrlString)
	query := requestUrl.Query()
//	query.Set("country",informationRequest.Location())
	query.Set("name",informationRequest.TopicName())

	fmt.Println("Forwarding Request to ",requestUrl)
	client := http.DefaultClient
	request,_ := http.NewRequest("GET",requestUrl.String(),nil)
	response, err := client.Do(request)
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