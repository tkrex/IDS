package scaling

import (
	"github.com/tkrex/IDS/common/models"
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/tkrex/IDS/domainController/persistence"
	"net/url"
)


// If the number of DomainInformationMessages for a Real World Domain,
// which is not the own, a new Domain Controller is requested at the Information Discovery Gateway
const ScalingThreshold = 10

//Manages Domain Controller Scaling Mechanism
type ScalingRequestManager struct {
	scalingInterfaceURL *url.URL
}

func NewScalingManager(scalingInterfaceURL *url.URL) *ScalingRequestManager {
	scalingManager := new(ScalingRequestManager)
	scalingManager.scalingInterfaceURL = scalingInterfaceURL
	return scalingManager
}

//Returns true if a new Domain Controller should be requested for a Real World Domain
func (scalingManager *ScalingRequestManager) CheckWorkloadForDomain(domain *models.RealWorldDomain) bool {
	dbWorker, error := persistence.NewDomainInformationStorage()
	if error != nil {
		fmt.Println(error)
		return false
	}
	dbWorker.Close()
	return dbWorker.NumberOfBrokersForDomain(domain) >= ScalingThreshold
}


//Creates a new Domain Controller by sending a request to the Information Discovery Gateway
func (scalingManager *ScalingRequestManager) CreateNewDominControllerForDomain(domain *models.RealWorldDomain) *models.DomainController {
	domainController, scalingError := scalingManager.requestNewDomainControllerForDomain(domain)
	if scalingError != nil {
		fmt.Println(scalingError)
		return nil
	}
	return domainController
}

//Requests a new Domain Controller at the Information Discovery Gateway
func (scalingManager *ScalingRequestManager) requestNewDomainControllerForDomain(domain *models.RealWorldDomain) (*models.DomainController, error) {
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
