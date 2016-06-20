package configuration

import "github.com/tkrex/IDS/common/models"

type DomainControllerConfigurationManager struct {

}


type DomainControllerConfiguration struct {
	DomainControllerID string `bson:"id"`
	ParentDomain *models.RealWorldDomain `bson:"parentDomain"`
}

func NewDomainControllerConfiguration(controllerID string,parentDomain *models.RealWorldDomain) *DomainControllerConfiguration {
	config := new(DomainControllerConfiguration)
	config.DomainControllerID = controllerID
	config.ParentDomain = parentDomain
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
		return err
	}
	defer storageManager.Close()
	conifg ,error := storageManager.FindDomainControllerConfig()
	return conifg, error
}
