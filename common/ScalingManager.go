package common

import (
	"github.com/tkrex/IDS/common/models"
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/tkrex/IDS/domainController/persistence"
	"net/url"
)

const ScalingThreshold = 10

type ScalingManager struct {
	scalingInterfaceURL *url.URL
}

func NewScalingManager(scalingInterfaceURL *url.URL) *ScalingManager {
	scalingManager := new(ScalingManager)
	scalingManager.scalingInterfaceURL = scalingInterfaceURL
	return scalingManager
}

func (scalingManager *ScalingManager) CheckWorkloadForDomain(domain *models.RealWorldDomain) bool {
	dbWorker, error := persistence.NewDomainControllerDatabaseWorker()
	if error != nil {
		fmt.Println(error)
		return false
	}
	dbWorker.Close()
	return dbWorker.NumberOfBrokersForDomain(domain) >= ScalingThreshold
}

func (scalingManager *ScalingManager) CreateNewDominControllerForDomain(domain *models.RealWorldDomain) *models.DomainController {
	domainController, scalingError := scalingManager.requestNewDomainControllerForDomain(domain)
	if scalingError != nil {
		fmt.Println(scalingError)
		return nil
	}
	return domainController
}

//TODO: Add Parent Domain as parameter
func (scalingManager *ScalingManager) requestNewDomainControllerForDomain(domain *models.RealWorldDomain) (*models.DomainController, error) {
	req, err := http.NewRequest("GET", scalingManager.scalingInterfaceURL.String() + "/rest/domainControllers/" + domain.Name + "/new", nil)
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