package providing

import (
	"github.com/tkrex/IDS/common/models"
	"github.com/tkrex/IDS/gateway/domainControllerManagement"
	"sync"
)

type RequestRoutingManager struct {
	dbManager *controlling.DomainControllerStorageManager
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
	routingManager.dbManager = controlling.NewDomainControllerStorageManager()
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
