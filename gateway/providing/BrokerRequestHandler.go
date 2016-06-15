package providing

import (
	"github.com/tkrex/IDS/common/controlling"
	"github.com/tkrex/IDS/common/models"

	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"errors"
	"net/url"
)

type BrokerRequestHandler struct {

}

func NewBrokerRequestHandler() *BrokerRequestHandler {
	return new(BrokerRequestHandler)
}


func (handler *BrokerRequestHandler) handleRequest(informationRequest *models.DomainInformationRequest) ([]*models.Broker,error) {


	//DEBUG CODE
	domain := models.NewRealWorldDomain("education")
	broker := models.NewBroker()
	broker.ID = "testID"
	broker.IP = "localhost"
	broker.RealWorldDomain = models.NewRealWorldDomain("education")
	broker.Geolocation = models.NewGeolocation("germany","bavaria","munich",11.6309,48.2499)
	broker.Statitics.NumberOfTopics = 20
	brokers := []*models.Broker{broker}
	return brokers, nil

	dbDelegate, err := controlling.NewControlMessageDBDelegate()
	if err != nil {
		return nil,err
	}
	defer dbDelegate.Close()

	var destinationDomainController *models.DomainController
	domain = models.NewRealWorldDomain(informationRequest.Domain())
	destinationDomainController = dbDelegate.FindDomainControllerForDomain(domain.FirstLevelDomain())

	if destinationDomainController == nil {
		destinationDomainController = dbDelegate.FindDomainControllerForDomain("default")
	}

	if destinationDomainController != nil {
		return handler.requestBrokersFromDomainController(informationRequest,destinationDomainController)
	}
	return nil, errors.New("No target controller found")
}

func (handler *BrokerRequestHandler) requestBrokersFromDomainController(informationRequest *models.DomainInformationRequest, domainController *models.DomainController) ([]*models.Broker,error) {
	requestUrlString := domainController.RestEndpoint.String() + "/brokers/" + informationRequest.Domain()
	requestUrl,_ := url.Parse(requestUrlString)
	query := requestUrl.Query()
	//query.Set("location",informationRequest.Location())
	query.Set("name",informationRequest.Name())

	fmt.Println("Forwarding Request to ",requestUrl)
	client := http.DefaultClient
	request,_ := http.NewRequest("GET",requestUrl.String(),nil)
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("%s", err)
		return nil,err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("%s", err)
		return nil,err
	}
	brokers := []*models.Broker{}
	json.Unmarshal([]byte(string(contents)), &brokers)
	return brokers,nil
}