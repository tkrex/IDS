package providing

import (
	"github.com/tkrex/IDS/common/models"
	"github.com/tkrex/IDS/common/controlling"
)

type ControllerForwardingManager struct {

}

func NewControllerForwardingManager() *ControllerForwardingManager {
	forwardManager := new(ControllerForwardingManager)
	return forwardManager
}

func (forwardingManager * ControllerForwardingManager) DomainControllerForDomain(domain *models.RealWorldDomain) *models.DomainController {
	dbDelegate, _ := controlling.NewControlMessageDBDelegate()
	if dbDelegate == nil {
		return nil
	}

	defer dbDelegate.Close()
	var destinationDomainController *models.DomainController
	destinationDomainController = dbDelegate.FindDomainControllerForDomain(domain.FirstLevelDomain())

	if destinationDomainController == nil {
		destinationDomainController = dbDelegate.FindDomainControllerForDomain("default")
	}
	return destinationDomainController
}