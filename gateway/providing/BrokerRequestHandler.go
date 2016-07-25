package providing

import (
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


func (handler *BrokerRequestHandler) handleRequest(informationRequest *models.DomainInformationRequest) (map[string][]*models.Broker,error) {

	domain := models.NewRealWorldDomain(informationRequest.Domain())
	targetDomain := domain.TopLevelDomain()
	destinationDomainController := RequestRoutingManagerInstance().DomainControllerForDomain(targetDomain)

	if destinationDomainController != nil {
		brokers, _ :=  handler.requestBrokersFromDomainController(informationRequest,destinationDomainController)
		if brokers != nil {
			sortedBrokers := handler.sortBrokersByDomains(brokers)
			return sortedBrokers, nil

		}
	}
	return nil, errors.New("No target controller found")
}


func (handler *BrokerRequestHandler) sortBrokersByDomains(brokers []*models.Broker) map[string][]*models.Broker {
	sortedBrokers := make(map[string][]*models.Broker)
	for _,broker := range brokers {
		domainBroker := sortedBrokers[broker.RealWorldDomain.Name]
		if domainBroker != nil {
			domainBroker = append(domainBroker,broker)
			sortedBrokers[broker.RealWorldDomain.Name] = domainBroker
		} else {
			sortedBrokers[broker.RealWorldDomain.Name] = []*models.Broker{broker}
		}
	}
	return sortedBrokers
}
func (handler *BrokerRequestHandler) requestBrokersFromDomainController(informationRequest *models.DomainInformationRequest, domainController *models.DomainController) ([]*models.Broker,error) {
	requestUrlString := domainController.RestEndpoint.String() + "/rest/brokers/" + informationRequest.Domain()
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