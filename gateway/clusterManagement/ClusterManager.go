package clusterManagement

import (
	"github.com/tkrex/IDS/common/models"
	"errors"
	"fmt"
)

type ClusterManager struct {
	clusterManagementStorage *ClusterManagerStorage
}

func NewClusterManager() *ClusterManager {
	worker := new(ClusterManager)
	worker.clusterManagementStorage = NewClusterManagerStorage()
	return worker
}

//Delegate Incoming ClusterManagementRequests to the corresponding method
func (manager *ClusterManager) HandleManagementRequest(request *models.ClusterManagementRequest) (*models.DomainController, error) {

	var requestError error
	var requestedDomainController  *models.DomainController

	switch request.RequestType {
	case models.DomainControllerStop:
		requestError = manager.stopDomainControllerInstance(request.Domain)
	case models.DomainControllerStart:
		requestedDomainController, requestError = manager.StartNewDomainControllerInstance(request.Domain, request.ParentDomain)

	case models.DomainControllerFetch:
		requestedDomainController, requestError = manager.domainControllerForDomain(request.Domain)
	}
	return requestedDomainController, requestError
}


//Returns a Domain Controller responsible for collecting information for the specified Real World Domain
func (manager *ClusterManager) domainControllerForDomain(domain *models.RealWorldDomain) (*models.DomainController, error) {
	var fetchError error
	var requestedDomainController *models.DomainController
	domainLevels := domain.DomainLevels()
	for i := len(domainLevels) - 1; i >= 0; i-- {
		fmt.Println("Searching Domain Controller for domain: ", domain)
		requestedDomainController,_ = manager.clusterManagementStorage.FindDomainControllerForDomain(domain)
		if requestedDomainController != nil {
			break
		}
		domain = domain.ParentDomain()
	}

	if requestedDomainController == nil {
		requestedDomainController,_ = manager.clusterManagementStorage.FindDomainControllerForDomain(models.NewRealWorldDomain("default"))
		if requestedDomainController == nil {
			fetchError = errors.New("No DomainController found")
		}
	}
	return requestedDomainController, fetchError
}

//Stops Domain Controller , which is assigned to the specifed Domain if existing
func (manager *ClusterManager) stopDomainControllerInstance(domain *models.RealWorldDomain) error {
	var stopError error
	existingDomainController,_ := manager.clusterManagementStorage.FindDomainControllerForDomain(domain)
	if existingDomainController != nil {
		stopError = NewDockerManager().StopDomainControllerInstance(domain)
		if stopError == nil {
			manager.clusterManagementStorage.RemoveDomainControllerForDomain(domain)
		} else {
			stopError = errors.New("Failed to start new domain controller instance")
		}
	} else {
		stopError = errors.New("No Domain Controller exists for this domain")
	}
	return stopError
}


//Starts a new Domain Controller for the specified Domain via the Docker Manager
func (manager *ClusterManager) StartNewDomainControllerInstance(domain *models.RealWorldDomain, parentDomain *models.RealWorldDomain) (*models.DomainController, error) {
	var startError error

	existingDomainController,_ := manager.clusterManagementStorage.FindDomainControllerForDomain(domain)
	if existingDomainController != nil {
		startError = errors.New("Domain Controller for this domain already exists")
	}
	domainController, startError := NewDockerManager().StartDomainControllerInstance(parentDomain, domain);
	if domainController != nil {
		manager.clusterManagementStorage.StoreDomainController(domainController)
	} else {
		startError = errors.New("Failed to start new domain controller instance")

	}
	return domainController, startError
}



