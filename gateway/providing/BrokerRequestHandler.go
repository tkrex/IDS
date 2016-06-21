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
	"github.com/tkrex/IDS/common/routing"
)

type BrokerRequestHandler struct {

}

func NewBrokerRequestHandler() *BrokerRequestHandler {
	return new(BrokerRequestHandler)
}


func (handler *BrokerRequestHandler) handleRequest(informationRequest *models.DomainInformationRequest) (map[string][]*models.Broker,error) {

	domains := []*models.RealWorldDomain{ models.NewRealWorldDomain("education/schools"), models.NewRealWorldDomain("education"), models.NewRealWorldDomain("education/university")}
	brokers := []*models.Broker{}
	for _, domain := range domains {
	//DEBUG CODE
	broker := models.NewBroker()
	broker.ID = "testID"
	broker.IP = "localhost"
	broker.RealWorldDomain = domain
	broker.Geolocation = models.NewGeolocation("germany", "bavaria", "munich", 11.6309, 48.2499)
	broker.Statitics.NumberOfTopics = 20
	brokers = append(brokers, broker)
	}
	sortedBrokers := handler.sortBrokersByDomains(brokers)
	return sortedBrokers, nil

	dbDelegate, err := controlling.NewControlMessageDBDelegate()
	if err != nil {
		return nil,err
	}
	defer dbDelegate.Close()

	domain := models.NewRealWorldDomain(informationRequest.Domain())
	targetDomain := domain.FirstLevelDomain()
	destinationDomainController,_ := routing.NewRoutingManager().DomainControllerForDomain(targetDomain,false)

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
		domains := sortedBrokers[broker.RealWorldDomain.Name]
		if domains != nil {
			domains = append(domains,broker)
		} else {
			sortedBrokers[broker.RealWorldDomain.Name] = []*models.Broker{broker}
		}
	}
	return sortedBrokers
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