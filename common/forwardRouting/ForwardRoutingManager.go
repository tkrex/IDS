package forwardRouting

import (
	"github.com/tkrex/IDS/common/models"
	"fmt"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"net/url"
	"time"
)

const RouteUpdateThreshold = time.Minute * 5

type ForwardRoutingManager struct {
	routingTable      map[string]*models.DomainController
	routeLifeTimes    map[string]time.Time
	clusterManagerURL *url.URL
}


func NewForwardRoutingManager(clusterManagerURL *url.URL) *ForwardRoutingManager {
	routingManager := new(ForwardRoutingManager)
	routingManager.routeLifeTimes = make(map[string]time.Time)
	routingManager.routingTable = make(map[string]*models.DomainController)
	routingManager.clusterManagerURL = clusterManagerURL
	return routingManager
}

func (routingManager *ForwardRoutingManager) DomainControllerForDomain(domain *models.RealWorldDomain, forceRefresh bool) (*models.DomainController, error) {
	if !forceRefresh {
		cachedDomainController, exist := routingManager.routingTable[domain.Name]
		if exist && time.Now().Sub(routingManager.routeLifeTimes[domain.Name]) <= RouteUpdateThreshold {
			return cachedDomainController, nil
		}
	}

	fetchedDomainController, err := routingManager.requestDomainControllerForDomain(domain)
	if fetchedDomainController != nil {
		routingManager.AddDomainControllerForDomain(fetchedDomainController,domain)
	}
	return fetchedDomainController, err
}

func (routingManager *ForwardRoutingManager) AddDomainControllerForDomain(domainController *models.DomainController, domain *models.RealWorldDomain) {
	routingManager.routingTable[domain.Name] = domainController
	routingManager.routeLifeTimes[domain.Name] = time.Now()
}

func (routingManager *ForwardRoutingManager) requestDomainControllerForDomain(domain *models.RealWorldDomain) (*models.DomainController, error) {
	fmt.Println("Sending Domain Controller Request for domain:",domain.Name)

	req, err := http.NewRequest("GET", routingManager.clusterManagerURL.String() + "/rest/domainControllers/" + domain.Name, nil)
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
