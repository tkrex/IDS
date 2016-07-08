package configuration

import (
	"github.com/tkrex/IDS/common/models"
	"net/url"
)

type DomainControllerConfigurationManager struct {

}


type DomainControllerConfiguration struct {
	DomainControllerID string `bson:"id"`
	ParentDomain *models.RealWorldDomain `bson:"parentDomain"`
	GatewayBrokerAddress *url.URL `bson:"gatewayBrokerAddress"`
}

func NewDomainControllerConfiguration(controllerID string,parentDomain *models.RealWorldDomain, gatewayBrokerAddress *url.URL) *DomainControllerConfiguration {
	config := new(DomainControllerConfiguration)
	config.DomainControllerID = controllerID
	config.ParentDomain = parentDomain
	config.GatewayBrokerAddress = gatewayBrokerAddress
	return config
}


func NewDomainControllerConfigurationManager() *DomainControllerConfigurationManager {
	configManager := new(DomainControllerConfigurationManager)
	return configManager
}

func (configManager *DomainControllerConfigurationManager) StoreConfig(config *DomainControllerConfiguration) error{
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
