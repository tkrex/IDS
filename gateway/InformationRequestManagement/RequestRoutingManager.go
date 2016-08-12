package informationRequestManagement

import (
	"github.com/tkrex/IDS/common/models"
	"sync"
	"github.com/tkrex/IDS/gateway/clusterManagement"
)

//Manages the routing of DomainInformationRequests to th Corresponding Top Level Domain Controller
type RequestRoutingManager struct {
	dbManager *clusterManagement.ClusterManagerStorage
}

var instance *RequestRoutingManager
var once sync.Once

func RequestRoutingManagerInstance() *RequestRoutingManager {
	once.Do(func() {
		instance = newRequestRoutingManager()
	})
	return instance
}


func newRequestRoutingManager() *RequestRoutingManager {
	routingManager := new(RequestRoutingManager)
	routingManager.dbManager = clusterManagement.NewClusterManagerStorage()
	return  routingManager
}


//Determines the TopLevelDomain of the specified Domain and fetches the Domain Controller from the ClusterManagementStorage
func (routingManager *RequestRoutingManager) DomainControllerForDomain(domain *models.RealWorldDomain) *models.DomainController {
	topLevelDomain := domain.TopLevelDomain()
	if domainController,_ := routingManager.dbManager.FindDomainControllerForDomain(topLevelDomain); domainController != nil {
		return domainController
	}
	if defaultDomainController,_ := routingManager.dbManager.FindDomainControllerForDomain(models.NewRealWorldDomain("default")); defaultDomainController != nil {
		return defaultDomainController
	}
	return nil
}
