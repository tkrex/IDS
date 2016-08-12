package configuration

import (
	"github.com/tkrex/IDS/common/models"
	"net/url"
	"os"
	"fmt"
	"sync"
)
//Singleton which manages the Domain Controller Configuration
type DomainControllerConfigurationManager struct {
	config *DomainControllerConfiguration
}

var instance *DomainControllerConfigurationManager
var once sync.Once

func DomainControllerConfigurationManagerInstance() *DomainControllerConfigurationManager {
	once.Do(func() {
		instance = NewDomainControllerConfigurationManager()
	})
	return instance
}


type DomainControllerConfiguration struct {
	DomainControllerID       string `bson:"id"`
	ParentDomain             *models.RealWorldDomain `bson:"parentDomain"`
	OwnDomain                *models.RealWorldDomain `bson:"ownDomain"`
	ControllerBrokerAddress  *url.URL `bson:"controllerBrokerAddress"`
	ClusterManagementAddress *url.URL `bson:"scalingInterfaceAddress"`
}

func NewDomainControllerConfiguration(controllerID string,parentDomain *models.RealWorldDomain, ownDomain *models.RealWorldDomain, controllerBrokerAddress *url.URL, clusterManagementAddress *url.URL) *DomainControllerConfiguration {
	config := new(DomainControllerConfiguration)
	config.DomainControllerID = controllerID
	config.ParentDomain = parentDomain
	config.OwnDomain = ownDomain
	config.ControllerBrokerAddress = controllerBrokerAddress
	config.ClusterManagementAddress = clusterManagementAddress
	return config
}


func NewDomainControllerConfigurationManager() *DomainControllerConfigurationManager {
	configManager := new(DomainControllerConfigurationManager)
	configManager.config = configManager.initConfig()
	return configManager
}
//Returns Domain Controller Config
func (configManager *DomainControllerConfigurationManager) Config() *DomainControllerConfiguration {
	return configManager.config
}

//Loads Domain Controller Configuration from Environment Variables
func (configManager *DomainControllerConfigurationManager) initConfig() *DomainControllerConfiguration {

	//Default values
	parentDomainString := "none"
	ownDomainString := "default"
	controllerID := "controllerID"
	brokerURLString := "ws://broker-default:9001"
	scalingInterfaceString := "http://localhost:8000"

	var parsingError error
	parentDomainString = os.Getenv("PARENT_DOMAIN")
	parentDomain := models.NewRealWorldDomain(parentDomainString)

	ownDomainString = os.Getenv("OWN_DOMAIN")
	ownDomain := models.NewRealWorldDomain(ownDomainString)

	controllerID = os.Getenv("CONTROLLER_ID")

	brokerURLString = os.Getenv("BROKER_URI")
	brokerURL,parsingError := url.Parse(brokerURLString)

	scalingInterfaceString = os.Getenv("MANAGEMENT_INTERFACE_URI")
	scalingInterfaceURL, parsingError := url.Parse(scalingInterfaceString)


	if parsingError != nil {
		fmt.Println("ConfigManager: Parsing Error")
	}

	config := NewDomainControllerConfiguration(controllerID, parentDomain,ownDomain, brokerURL,scalingInterfaceURL)
	return config
}
