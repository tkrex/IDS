package configuration

import (
	"net/url"
	"os"
	"fmt"
	"sync"
)

type DaemonConfigurationManager struct {
	config *DaemonConfiguration
}

var instance *DaemonConfigurationManager
var once sync.Once

func DaemonConfigurationManagerInstance() *DaemonConfigurationManager {
	once.Do(func() {
		instance = newDaemonConfigurationManager()
	})
	return instance
}

type DaemonConfiguration struct {
	BrokerURL            *url.URL `bson:"controllerBrokerAddress"`
	RoutingManagementURL *url.URL `bson:"routingManagementAddress"`
	RegistrationURL      *url.URL `bson:"registrationURL"`
	DatabaseAddress      string `bson:"databaseAddress"`
}

func NewDaemonConfiguration(brokerURL *url.URL, routingManagementURL *url.URL,registrationURL *url.URL, databaseURL string) *DaemonConfiguration {
	config := new(DaemonConfiguration)
	config.BrokerURL = brokerURL
	config.RoutingManagementURL = routingManagementURL
	config.DatabaseAddress = databaseURL
	return config
}

func newDaemonConfigurationManager() *DaemonConfigurationManager {
	configManager := new(DaemonConfigurationManager)
	configManager.config = configManager.initConfig()
	fmt.Println(configManager.config)
	return configManager
}

func (configManager *DaemonConfigurationManager) Config() *DaemonConfiguration {
	return configManager.config
}

func (configManager *DaemonConfigurationManager) initConfig() *DaemonConfiguration {


	brokerURLString := "ws://localhost:18833"
	registrationURLString := "http://localhost:8000"
	mongoURL := "localhost:27017"
	routingManagementURL,_ := url.Parse("http://localhost:8000")

	if existingConfig, _ := configManager.fetchDomainControllerConfig(); existingConfig != nil {
		configManager.config = existingConfig
		return existingConfig
	}


	mongoURL = os.Getenv("MONGODB_URI")

	brokerURLString = os.Getenv("BROKER_URI")
	brokerURL, parsingError := url.Parse(brokerURLString)

	registrationURLString = os.Getenv("REGISTRATION_URL")
	registrationURL, parsingError := url.Parse(registrationURLString)

	if parsingError != nil {
		fmt.Println("ConfigManager: Parsing Error")
	}

	config := NewDaemonConfiguration(brokerURL,routingManagementURL,registrationURL,mongoURL)
	configManager.storeConfig(config)
	return config
}

func (configManager *DaemonConfigurationManager) storeConfig(config *DaemonConfiguration) error {
	storageManager, err := NewDaemonConfigStorageManager()
	if err != nil {
		return err
	}
	err = storageManager.StoreDaemonConfig(config)
	return err
}

func (configManager *DaemonConfigurationManager) fetchDomainControllerConfig() (*DaemonConfiguration, error) {
	storageManager, err := NewDaemonConfigStorageManager()
	if err != nil {
		return nil, err
	}
	defer storageManager.Close()
	conifg, error := storageManager.FindDaemonConfig()
	return conifg, error
}

