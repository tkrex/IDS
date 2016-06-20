package routing

import (
	"github.com/tkrex/IDS/common/models"
	"fmt"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"net/url"
	"time"
	"github.com/tkrex/IDS/common/controlling"
)

type RoutingManager struct {
	dbDelegate *controlling.ControlMessageDBDelegate
	lastUpdate 	time.Time
}


func NewRoutingManager() *RoutingManager {
	routingManager := new(RoutingManager)
	dbDelegate,_ := controlling.NewControlMessageDBDelegate()
	routingManager.dbDelegate = dbDelegate

	return routingManager
}

func (routingManager *RoutingManager) DomainControllerForDomain(domain *models.RealWorldDomain, forceRefresh bool) (*models.DomainController,error) {

	if !forceRefresh {
		cachedDomainController := routingManager.dbDelegate.FindDomainControllerForDomain(domain)
		if cachedDomainController != nil {
			return cachedDomainController, nil
		}
	}
	fetchedDomainController, err := routingManager.requestDomainControllerForDomain(domain)
	if fetchedDomainController != nil {
		routingManager.AddDomainController(fetchedDomainController)
	}
	return fetchedDomainController, err
}

func (routingManager *RoutingManager) AddDomainController(domainController *models.DomainController) {
	routingManager.dbDelegate.StoreDomainController(domainController)
}



func (routingManager *RoutingManager) requestDomainControllerForDomain(domain *models.RealWorldDomain) (*models.DomainController,error) {
	fmt.Println("Sending Broker Registration Request")

	InfrastructureManagerURL,_ := url.Parse("http://localhost:8080/rest")

	req, err := http.NewRequest("GET",  InfrastructureManagerURL.String() + "/domainControllers/"+domain.Name,nil)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var domainController *models.DomainController
	err = json.Unmarshal(body, &domainController)
	return domainController, err
}





