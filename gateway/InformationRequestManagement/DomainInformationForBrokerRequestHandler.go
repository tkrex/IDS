package informationRequestManagement

import (
	"net/http"
	"github.com/tkrex/IDS/common/models"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"net/url"
	"errors"
)

type DomainInformationForBrokerRequestHandler struct {

}

func NewDomainInformationForBrokerRequestHandler() *DomainInformationForBrokerRequestHandler {
	return new(DomainInformationForBrokerRequestHandler)
}


func (handler *DomainInformationForBrokerRequestHandler) HandleRequest(brokerId string, informationRequest *models.DomainInformationRequest) (*models.DomainInformationMessage,error) {


	destinationDomainController := RequestRoutingManagerInstance().DomainControllerForDomain(models.NewRealWorldDomain(informationRequest.Domain()))
	if destinationDomainController == nil {
		return nil, errors.New("No target controller found")
	}
	return handler.forwardRequestToDomainController(brokerId,informationRequest,destinationDomainController)
}
//Forwards the request to the corresponding Top Level Domain Controller and returns the results
func (handler *DomainInformationForBrokerRequestHandler) forwardRequestToDomainController(brokerId string, informationRequest *models.DomainInformationRequest,domainController *models.DomainController) (*models.DomainInformationMessage,error) {
	requestUrlString := domainController.RestEndpoint.String() + "/rest/brokers/" + brokerId + "/" + informationRequest.Domain()
	requestUrl,_ := url.Parse(requestUrlString)
	query := requestUrl.Query()

	locationJSON, _ := json.Marshal(informationRequest.Location())

        query.Set("location",string(locationJSON))
	query.Set("name",informationRequest.TopicName())

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