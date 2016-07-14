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

func DaemonConfigurationManager() *DaemonConfigurationManager {
	once.Do(func() {
		instance = NewDaemonConfigurationManager()
	})
	return instance
}

type DaemonConfiguration struct {
	BrokerURL            *url.URL `bson:"controllerBrokerAddress"`
	RoutingManagementURL *url.URL `bson:"routingManagementAddress"`
	DatabaseURL *url.URL `bson:"databaseURL"`
}

func NewDaemonConfiguration(brokerURL *url.URL, routingManagementURL *url.URL, databaseURL *url.URL) *DaemonConfiguration {
	config := new(DaemonConfiguration)
	config.BrokerURL = brokerURL
	config.RoutingManagementURL = routingManagementURL
	config.DatabaseURL = databaseURL
	return config
}

func NewDaemonConfigurationManager() *DaemonConfigurationManager {
	configManager := new(DaemonConfigurationManager)
	configManager.config = configManager.initConfig()
	return configManager
}

func (configManager *DaemonConfigurationManager) Config() *DaemonConfiguration {
	return configManager.config
}

func (configManager *DaemonConfigurationManager) initConfig() *DaemonConfiguration {


	brokerURLString := "ws://localhost:18833"
	mongoURLString := "localhost:27017"
	routingManagementURL,_ := url.Parse("http://localhost:8000")

	if existingConfig, _ := configManager.fetchDomainControllerConfig(); existingConfig != nil {
		configManager.config = existingConfig
		return existingConfig
	}


	mongoURLString = os.Getenv("MONGODB_URI")
	mongoURL, parsingError := url.Parse(mongoURLString)
	brokerURLString = os.Getenv("BROKER_URI")
	brokerURL, parsingError := url.Parse(brokerURLString)


	if parsingError != nil {
		fmt.Println("ConfigManager: Parsing Error")
	}

	config := NewDaemonConfiguration(brokerURL,routingManagementURL,mongoURL)
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

