package configuration

import (
	"github.com/tkrex/IDS/common/models"
	"net/url"
	"os"
	"fmt"
	"github.com/tkrex/IDS/daemon/configuration"
)

type DomainControllerConfigurationManager struct {

}


type DomainControllerConfiguration struct {
	DomainControllerID string `bson:"id"`
	ParentDomain *models.RealWorldDomain `bson:"parentDomain"`
	OwnDomain *models.RealWorldDomain `bson:"ownDomain"`
	ControllerBrokerAddress *url.URL `bson:"controllerBrokerAddress"`
	GatewayBrokerAddress *url.URL `bson:"gatewayBrokerAddress"`

}

func NewDomainControllerConfiguration(controllerID string,parentDomain *models.RealWorldDomain, ownDomain *models.RealWorldDomain, controllerBrokerAddress *url.URL, gatewayBrokerAddress *url.URL) *DomainControllerConfiguration {
	config := new(DomainControllerConfiguration)
	config.DomainControllerID = controllerID
	config.ParentDomain = parentDomain
	config.OwnDomain = ownDomain
	config.ControllerBrokerAddress = controllerBrokerAddress
	config.GatewayBrokerAddress = gatewayBrokerAddress
	return config
}


func NewDomainControllerConfigurationManager() *DomainControllerConfigurationManager {
	configManager := new(DomainControllerConfigurationManager)
	return configManager
}

func (configManager *DomainControllerConfigurationManager) InitConfig() *DomainControllerConfiguration {
	parentDomainString := os.Getenv("PARENT_DOMAIN")
	if parentDomainString == "" {
		parentDomainString = "default"
	}

	parentDomain := models.NewRealWorldDomain(parentDomainString)

	ownDomainString := os.Getenv("OWN_DOMAIN")
	if parentDomainString == "" {
		parentDomainString = "default"
	}

	ownDomain := models.NewRealWorldDomain(ownDomainString)


	controllerID := os.Getenv("CONTROLLER_ID")
	if controllerID == "" {
		controllerID = "controllerID"
	}

	brokerURLString := os.Getenv("BROKER_URI")
	fmt.Println()
	if brokerURLString == "" {
		brokerURLString = "ws://localhost:18833"
	}
	brokerURL,error := url.Parse(brokerURLString)
	if error != nil {
		fmt.Println("Parsing Error: ",error)
	}

	gatewayBrokerURLString := os.Getenv("GATEWAY_BROKER_URI")
	if gatewayBrokerURLString == "" {
		gatewayBrokerURLString = "ws://localhost:18833"
	}
	gatewayBrokerURL,_ := url.Parse(gatewayBrokerURLString)
	config := NewDomainControllerConfiguration(controllerID, parentDomain,ownDomain, brokerURL,gatewayBrokerURL)
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

func (configManager *DomainControllerConfigurationManager) DomainControllerConfig() (*DomainControllerConfiguration,error) {
	storageManager, err := NewDomainControllerConfigStorageManager()
	if err != nil {
		return nil, err
	}
	defer storageManager.Close()
	conifg ,error := storageManager.FindDomainControllerConfig()
	return conifg, error
}
