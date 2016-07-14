package controlling

import (
	"github.com/tkrex/IDS/common/models"
	"github.com/tkrex/IDS/common/controlling"
	"github.com/tkrex/IDS/gateway/scaling"
	"errors"
	"fmt"
)

type DomainControllerManager struct {

}

const defaultDomain = models.NewRealWorldDomain("default")

func NewDomainControllerManager() *DomainControllerManager {
	worker := new(DomainControllerManager)
	return worker
}

func (handler *DomainControllerManager) handleManagementRequest(request *models.DomainControllerManagementRequest) (*models.DomainController, error) {

	var requestError error
	var requestedDomainController  *models.DomainController

	switch request.RequestType {
	case models.DomainControllerStop:
		requestError = handler.stopDomainControllerInstance(request.Domain)
	case models.DomainControllerStart:
		requestedDomainController, requestError = handler.startNewDomainControllerInstance(request.Domain, request.ParentDomain)

	case models.DomainControllerFetch:
		requestedDomainController , requestError = handler.DomainControllerForDomain(request.Domain)
	}
	return requestedDomainController, requestError
}

func (handler *DomainControllerManager) DomainControllerForDomain(domain *models.RealWorldDomain) (*models.DomainController, error) {
	var fetchError error
	var requestedDomainController *models.DomainController
	dbWorker, fetchError := controlling.NewControlMessageDBDelegate()
	domainLevels := domain.DomainLevels()
	for i := len(domainLevels)-1; i >= 0; i-- {
		fmt.Println("Searching Domain Controller for domain: ",domain)
		requestedDomainController= dbWorker.FindDomainControllerForDomain(domain)
		if requestedDomainController != nil {
			break
		}
		domain = domain.ParentDomain()
	}

	if requestedDomainController == nil {
		requestedDomainController= dbWorker.FindDomainControllerForDomain(domain)
		if requestedDomainController == nil {
			fetchError = errors.New("No DomainController found")
		}
	}
	return requestedDomainController, fetchError
}

func (handler *DomainControllerManager) stopDomainControllerInstance(domain *models.RealWorldDomain) error {
	var stopError error
	dbWorker, error := controlling.NewControlMessageDBDelegate()
	if error != nil {
		stopError = errors.New("Failed to start new domain controller instance")
	}
	defer dbWorker.Close()
	existingDomainController := dbWorker.FindDomainControllerForDomain(domain)
	if existingDomainController != nil {
		stopError = scaling.NewDockerManager().StopDomainControllerInstance(domain)
		if stopError == nil {
			dbWorker.RemoveDomainControllerForDomain(domain)
		} else {
			stopError = errors.New("Failed to start new domain controller instance")
		}
	} else {
		stopError = errors.New("No Domain Controller exists for this domain")
	}
	return  stopError
}


func (handler *DomainControllerManager) startNewDomainControllerInstance(domain *models.RealWorldDomain, parentDomain *models.RealWorldDomain) (*models.DomainController, error) {
	var startError error
	dbWorker, error := controlling.NewControlMessageDBDelegate()
	if error != nil {
		startError = errors.New("Failed to start new domain controller instance")
	}
	defer dbWorker.Close()

	existingDomainController := dbWorker.FindDomainControllerForDomain(domain)
	if existingDomainController != nil {
		startError = errors.New("Domain Controller for this domain already exists")
	}
	domainController, error := scaling.NewDockerManager().StartDomainControllerInstance(parentDomain, domain)
	if error != nil {
		startError = errors.New("Failed to start new domain controller instance")
	}
	return domainController, startError
}



