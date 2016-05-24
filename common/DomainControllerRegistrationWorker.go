package common

import (
	"github.com/tkrex/IDS/common/models"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

type DomainControllerRegistrationWorker struct {

}


func NewDomainControllerRegistrationWorker() *DomainControllerRegistrationWorker {
	return new(DomainControllerRegistrationWorker)
}

func (worker *DomainControllerRegistrationWorker) RequestNewDomainControllerForDomain(domain *models.RealWorldDomain) {
	domainController , err := worker.sendManagementRequest(models.DomainControllerChange,domain)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(domainController)
	}
}

func (worker *DomainControllerRegistrationWorker) RequestDomainControllerDeletionForDomain(domain *models.RealWorldDomain) {
	domainController , err := worker.sendManagementRequest(models.DomainControllerDelete,domain)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(domainController)
	}
}

func  (worker *DomainControllerRegistrationWorker) sendManagementRequest(requestType models.ControlMessageType,domain *models.RealWorldDomain) (*models.DomainController, error) {
	//TODO: Get Gateway Address

	registrationURL := "http://localhost:8000/domainController/"+domain.Name

	switch requestType {
	case models.DomainControllerChange:
		registrationURL += "/new"
	case models.DomainControllerDelete:
		registrationURL += "/delete"

	}
	fmt.Println("Sending Management Request: " + registrationURL)
	var domainController *models.DomainController
	resp, err := http.Get(registrationURL)
	if err != nil {

		return nil, err
	}
	defer resp.Body.Close()
	fmt.Println("Management Request: " + resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body,&domainController)
	if err != nil {
		return nil ,err
	}
	return domainController, nil
}