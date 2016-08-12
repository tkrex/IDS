package configuration

import (
	"net/url"
	"os"
	"fmt"
	"sync"
)

//Singleton, which initilazing the config of the daemon via Environment Variables
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
	config.RegistrationURL = registrationURL
	return config
}

func newDaemonConfigurationManager() *DaemonConfigurationManager {
	configManager := new(DaemonConfigurationManager)
	configManager.config = configManager.initConfig()
	return configManager
}
//Returns Configuration
func (configManager *DaemonConfigurationManager) Config() *DaemonConfiguration {
	return configManager.config
}

//Creates Configuration by loading Environment Variables of Host OS
func (configManager *DaemonConfigurationManager) initConfig() *DaemonConfiguration {
	fmt.Println("Init Config")

	brokerURLString := "tcp://localhost:1883"
	registrationURLString := "http://localhost:8000"
	mongoURL := "localhost:27017"
	routingManagementURL,_ := url.Parse("http://localhost:8000")

	mongoURL = os.Getenv("MONGODB_URI")

	if broker := os.Getenv("BROKER_URI"); broker != "" {
		brokerURLString = broker
	}
	brokerURL, parsingError := url.Parse(brokerURLString)

	if registration  := os.Getenv("REGISTRATION_URL"); registration != "" {
		registrationURLString = registration
	}
	registrationURL, parsingError := url.Parse(registrationURLString)

	if parsingError != nil {
		fmt.Println("ConfigManager: Parsing Error")
	}

	config := NewDaemonConfiguration(brokerURL,routingManagementURL,registrationURL,mongoURL)
	return config
}

