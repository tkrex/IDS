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
	"sync"
	"github.com/tkrex/IDS/domainController/configuration"
)

const RouteUpdateThreshold = time.Minute * 5

type RoutingManager struct {
	routingTable   map[string]*models.DomainController
	routeLifeTimes map[string]time.Time
}


func NewRoutingManager() *RoutingManager {
	routingManager := new(RoutingManager)
	routingManager.routeLifeTimes = make(map[string]time.Time)
	routingManager.routingTable = make(map[string]*models.DomainController)
	return routingManager
}

var instance *RoutingManager
var once sync.Once

func RoutingManager() *RoutingManager {
	once.Do(func() {
		instance = NewRoutingManager()
	})
	return instance
}


func (routingManager *RoutingManager) DomainControllerForDomain(domain *models.RealWorldDomain, forceRefresh bool) (*models.DomainController, error) {
	if !forceRefresh {
		cachedDomainController, exist := routingManager.routingTable[domain.Name]
		if exist && time.Now().Sub(routingManager.routeLifeTimes[domain.Name]).Minutes() <= RouteUpdateThreshold {
			return cachedDomainController, nil
		}
	}

	fetchedDomainController, err := routingManager.requestDomainControllerForDomain(domain)
	if fetchedDomainController != nil {
		routingManager.AddDomainControllerForDomain(fetchedDomainController,domain)
	}
	return fetchedDomainController, err
}

func (routingManager *RoutingManager) AddDomainControllerForDomain(domainController *models.DomainController, domain *models.RealWorldDomain) {
	routingManager.routingTable[domain.Name] = domainController
	routingManager.routeLifeTimes[domain.Name] = time.Now()
}

func (routingManager *RoutingManager) requestDomainControllerForDomain(domain *models.RealWorldDomain) (*models.DomainController, error) {
	fmt.Println("Sending Domain Controller Request for domain:",domain.Name)
	infrastructureManagerURL := configuration.DomainControllerConfigurationManager().Config().ScalingInterfaceAddress

	req, err := http.NewRequest("GET", infrastructureManagerURL.String() + "/rest/domainControllers/" + domain.Name, nil)
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
