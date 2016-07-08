package main

import (

	"github.com/tkrex/IDS/common/models"
	"github.com/tkrex/IDS/common/subscribing"
	"github.com/tkrex/IDS/domainController/processing"
	"github.com/tkrex/IDS/domainController/providing"
	"os"
	"fmt"
	"net/url"
	"github.com/tkrex/IDS/domainController/configuration"
)



func main() {

	initDomainControllerConfiguration()
	go startDomainInformationProcessing()
	_ = providing.NewDomainInformationRESTProvider("8080")
	for {}
}


func startDomainInformationProcessing() {
	//producer layer

	brokerAddressString := os.Getenv("BROKER_URI")
	if brokerAddressString == "" {
		brokerAddressString = "ws://localhost:11883"
	}
	brokerAddress, _ := url.Parse(brokerAddressString)
	fmt.Println("Broker URI: ",brokerAddress)

	desiredTopic  := "#"
	subscriberConfig := models.NewMqttClientConfiguration(brokerAddress,"domainController")
	subscriber := subscribing.NewMqttSubscriber(subscriberConfig,desiredTopic,false)
	//processing layer
	_ = processing.NewDomainInformationProcessor(subscriber.IncomingTopicsChannel())
}

func initDomainControllerConfiguration () {
	parentDomainString := os.Getenv("PARENT_DOMAIN")
	if parentDomainString == "" {
		parentDomainString = "default"
	}

	parentDomain := models.NewRealWorldDomain(parentDomainString)


	controllerID := os.Getenv("CONTROLLER_ID")
	if controllerID == "" {
		controllerID = "controllerID"
	}

	brokerURLString := os.Getenv("BROKER_URL")
	if brokerURLString == "" {
		brokerURLString = "ws://localhost:18833"
	}
	brokerURL,_ := url.Parse(brokerURLString)

	config := configuration.NewDomainControllerConfiguration(controllerID, parentDomain, brokerURL)
	configuration.NewDomainControllerConfigurationManager().StoreConfig(config)
}
