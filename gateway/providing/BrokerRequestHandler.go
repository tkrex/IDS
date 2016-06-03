package providing

import (
	"github.com/tkrex/IDS/common/controlling"
	"github.com/tkrex/IDS/common/models"

)

type BrokerRequestHandler struct {

}

func NewBrokerRequestHandler() *BrokerRequestHandler {
	return new(BrokerRequestHandler)
}


func (handler *BrokerRequestHandler) handleRequest(domainName string) []*models.Broker {


	//DEBUG CODE
	domain := models.NewRealWorldDomain("education")
	broker := models.NewBroker()
	broker.ID = "testID"
	broker.IP = "localhost"
	broker.RealWorldDomain = models.NewRealWorldDomain("education")
	broker.Geolocation = models.NewGeolocation("germany","bavaria","munich",11.6309,48.2499)
	broker.Statitics.NumberOfTopics = 20
	brokers := []*models.Broker{broker}
	return brokers

	dbDelegate, err := controlling.NewControlMessageDBDelegate()
	if err != nil {
		return nil
	}
	defer dbDelegate.Close()

	var destinationDomainController *models.DomainController
	domain = models.NewRealWorldDomain(domainName)
	destinationDomainController = dbDelegate.FindDomainControllerForDomain(domain.FirstLevelDomain())

	if destinationDomainController == nil {
		destinationDomainController = dbDelegate.FindDomainControllerForDomain("default")
	}

	if destinationDomainController != nil {
		return handler.requestBrokersFromDomainController(domainName,destinationDomainController)
	}
	return nil
}


//TODO: ADD Broker Request to DOmain COntroller
func (handler *BrokerRequestHandler) requestBrokersFromDomainController(domainName string, domainController *models.DomainController) []*models.Broker {
	//requestUrl := "http://" + domainController.IpAddress + ":8080/domainController/domainInformation/" + domainName
	//fmt.Println("Forwarding Request to ",requestUrl)
	//client := http.DefaultClient
	//response, err := client.Get(requestUrl)
	//if err != nil {
	//	fmt.Printf("%s", err)
	//	return nil
	//}
	//defer response.Body.Close()
	//contents, err := ioutil.ReadAll(response.Body)
	//if err != nil {
	//	fmt.Printf("%s", err)
	//	return nil
	//}
	//domainInformation := []*models.DomainInformationMessage{}
	//json.Unmarshal([]byte(string(contents)), &domainInformation)
	//return domainInformation
	return []*models.Broker{}
}