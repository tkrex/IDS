package providing

import (
	"github.com/tkrex/IDS/common/models"
	"sync"
	"github.com/tkrex/IDS/gateway/domainControllerManagement"
)

type RequestRoutingManager struct {
	dbManager *domainControllerManagement.DomainControllerStorageManager
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
	routingManager.dbManager = domainControllerManagement.NewDomainControllerStorageManager()
	return  routingManager
}

func (routingManager *RequestRoutingManager) DomainControllerForDomain(domain *models.RealWorldDomain) *models.DomainController {
	topLevelDomain := domain.TopLevelDomain()
	if domainController := routingManager.dbManager.FindDomainControllerForDomain(topLevelDomain); domainController != nil {
		return domainController
	}
	if defaultDomainController := routingManager.dbManager.FindDomainControllerForDomain(models.NewRealWorldDomain("default")); defaultDomainController != nil {
		return defaultDomainController
	}
	return nil
}
