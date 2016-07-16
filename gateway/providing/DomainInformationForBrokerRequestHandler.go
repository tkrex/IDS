package providing

import (
	"net/http"
	"github.com/tkrex/IDS/common/models"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"net/url"
	"errors"
	"time"
)

type DomainInformationForBrokerRequestHandler struct {

}

func NewDomainInformationForBrokerRequestHandler() *DomainInformationForBrokerRequestHandler {
	return new(DomainInformationForBrokerRequestHandler)
}


func (handler *DomainInformationForBrokerRequestHandler) handleRequest(brokerId string, informationRequest *models.DomainInformationRequest) (*models.DomainInformationMessage,error) {
	domain := models.NewRealWorldDomain(informationRequest.Domain())

	//DEBUG CODE
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
	return message,nil

	destinationDomainController := RequestRoutingManagerInstance().DomainControllerForDomain(domain)
	if destinationDomainController == nil {
		return nil, errors.New("No target controller found")
	}
	return handler.forwardRequestToDomainController(brokerId,informationRequest,destinationDomainController)
}

func (handler *DomainInformationForBrokerRequestHandler) forwardRequestToDomainController(brokerId string, informationRequest *models.DomainInformationRequest,domainController *models.DomainController) (*models.DomainInformationMessage,error) {
	requestUrlString := domainController.RestEndpoint.String() + "/brokers/" + brokerId + "/" + informationRequest.Domain()
	requestUrl,_ := url.Parse(requestUrlString)
	query := requestUrl.Query()

	locationJSON, _ := json.Marshal(informationRequest.Location())

        query.Set("location",string(locationJSON))
	query.Set("name",informationRequest.Name())

	fmt.Println("Forwarding Request to ",requestUrl)
	client := http.DefaultClient
	request,_ := http.NewRequest("GET",requestUrl.String(),nil)
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var domainInformation *models.DomainInformationMessage
	json.Unmarshal([]byte(string(contents)), &domainInformation)
	return domainInformation, nil
}