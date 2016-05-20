package providing

import (
	"net/http"
	"github.com/tkrex/IDS/common/controlling"
	"github.com/tkrex/IDS/common/models"
	"fmt"
	"io/ioutil"
	"encoding/json"
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
	broker := models.NewBroker("localhost","krex.com")
	topics := []*models.Topic{}
	message := models.NewDomainInformationMessage(domain,broker,topics)
	domainInformation = append(domainInformation,message)
	return domainInformation

	dbDelegate, err := controlling.NewControlMessageDBDelegate()
	if err != nil {
		return nil
	}
	defer dbDelegate.Close()

	var destinationDomainController *models.DomainController
	destinationDomainController = dbDelegate.FindDomainControllerForDomain(domainName)

	if destinationDomainController == nil {
		destinationDomainController = dbDelegate.FindDomainControllerForDomain("default")
	}

	if destinationDomainController != nil {
		return handler.requestDomainInformationFromDomainController(domainName,destinationDomainController)
	}
	return nil
}


func (handler *DomainInformationRequestHandler) requestDomainInformationFromDomainController(domainName string, domainController *models.DomainController) []*models.DomainInformationMessage {
	requestUrl := "http://" + domainController.IpAddress + "/domainInformation/" + domainName
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