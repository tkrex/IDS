package informationRequestManagement

import (
	"github.com/tkrex/IDS/common/models"
	"sync"
	"github.com/tkrex/IDS/gateway/clusterManagement"
)

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
