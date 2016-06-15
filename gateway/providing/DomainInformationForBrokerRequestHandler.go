package providing

import (
	"net/http"
	"github.com/tkrex/IDS/common/controlling"
	"github.com/tkrex/IDS/common/models"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"time"
	"github.com/tkrex/IDS/common/routing"
)

type DomainInformationForBrokerRequestHandler struct {

}

func NewDomainInformationForBrokerRequestHandler() *DomainInformationForBrokerRequestHandler {
	return new(DomainInformationForBrokerRequestHandler)
}


func (handler *DomainInformationForBrokerRequestHandler) handleRequest(brokerId string, informationRequest *models.DomainInformationRequest) (*models.DomainInformationMessage,error) {
	destinationDomainController := routing.NewRoutingManager().DomainControllerForDomain(informationRequest.Domain())
	if destinationDomainController == nil {
		return nil, error("No target controller found")
	}
	return handler.forwardRequestToDomainController(brokerId,informationRequest,destinationDomainController)
}

func (handler *DomainInformationForBrokerRequestHandler) forwardRequestToDomainController(brokerId string, informationRequest *models.DomainInformationRequest,domainController *models.DomainController) (*models.DomainInformationMessage,error) {
	requestUrl := domainController.RestEndpoint.String() + "/brokers/" + brokerId + "/" + informationRequest.Domain()
	fmt.Println("Forwarding Request to ",requestUrl)
	client := http.DefaultClient
	request,_ := http.NewRequest("GET",requestUrl,nil)
	request.FormValue("country") = informationRequest.Country()
	request.FormValue("name") = informationRequest.Name()
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("%s", err)
		return nil, err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("%s", err)
		return nil, err
	}
	var domainInformation *models.DomainInformationMessage
	json.Unmarshal([]byte(string(contents)), &domainInformation)
	return domainInformation, nil
}