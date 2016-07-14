package configuration

import (
	"github.com/tkrex/IDS/common/models"
	"net/url"
	"os"
	"fmt"
	"github.com/tkrex/IDS/daemon/configuration"
	"sync"
)

type DomainControllerConfigurationManager struct {
	config *DomainControllerConfiguration
}

var instance *DomainControllerConfigurationManager
var once sync.Once

func DomainControllerConfigurationManager() *DomainControllerConfigurationManager {
	once.Do(func() {
		instance = NewDomainControllerConfigurationManager()
	})
	return instance
}


type DomainControllerConfiguration struct {
	DomainControllerID string `bson:"id"`
	ParentDomain *models.RealWorldDomain `bson:"parentDomain"`
	OwnDomain *models.RealWorldDomain `bson:"ownDomain"`
	ControllerBrokerAddress *url.URL `bson:"controllerBrokerAddress"`
	GatewayBrokerAddress *url.URL `bson:"gatewayBrokerAddress"`
	ScalingInterfaceAddress *url.URL `bson:"scalingInterfaceAddress"`
}

func NewDomainControllerConfiguration(controllerID string,parentDomain *models.RealWorldDomain, ownDomain *models.RealWorldDomain, controllerBrokerAddress *url.URL, gatewayBrokerAddress *url.URL, scalingInterfaceAddress *url.URL) *DomainControllerConfiguration {
	config := new(DomainControllerConfiguration)
	config.DomainControllerID = controllerID
	config.ParentDomain = parentDomain
	config.OwnDomain = ownDomain
	config.ControllerBrokerAddress = controllerBrokerAddress
	config.GatewayBrokerAddress = gatewayBrokerAddress
	config.ScalingInterfaceAddress = scalingInterfaceAddress
	return config
}


func NewDomainControllerConfigurationManager() *DomainControllerConfigurationManager {
	configManager := new(DomainControllerConfigurationManager)
	configManager.config = configManager.initConfig()
	return configManager
}

func (configManager *DomainControllerConfigurationManager) Config() *DomainControllerConfiguration {
	return configManager.config
}

func (configManager *DomainControllerConfigurationManager) initConfig() *DomainControllerConfiguration {

	//Default values
	parentDomainString := "default"
	ownDomainString := "default"
	controllerID := "controllerID"
	brokerURLString := "ws://localhost:18833"
	gatewayBrokerURLString := "ws://localhost:18833"
	scalingInterfaceString := "http://localhost:8000"


	if existingConfig, _ :=  configManager.fetchDomainControllerConfig(); existingConfig != nil {
		configManager.config = existingConfig
		return existingConfig
	}


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

	gatewayBrokerURLString = os.Getenv("GATEWAY_BROKER_URI")
	gatewayBrokerURL,parsingError := url.Parse(gatewayBrokerURLString)

	if parsingError != nil {
		fmt.Println("ConfigManager: Parsing Error")
	}

	config := NewDomainControllerConfiguration(controllerID, parentDomain,ownDomain, brokerURL,gatewayBrokerURL,scalingInterfaceURL)
	configManager.storeConfig(config)
	return config
}


func (configManager *DomainControllerConfigurationManager) storeConfig(config *DomainControllerConfiguration) error{
	storageManager, err := NewDomainControllerConfigStorageManager()
	if err != nil {
		return err
	}
	err = storageManager.StoreDomainControllerConfig(config)
	return err
}

func (configManager *DomainControllerConfigurationManager) fetchDomainControllerConfig() (*DomainControllerConfiguration,error) {
	storageManager, err := NewDomainControllerConfigStorageManager()
	if err != nil {
		return nil, err
	}
	defer storageManager.Close()
	conifg ,error := storageManager.FindDomainControllerConfig()
	return conifg, error
}
