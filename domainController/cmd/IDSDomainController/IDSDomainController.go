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


	config, _ := configuration.NewDomainControllerConfigurationManager().DomainControllerConfig()

	fmt.Println("Broker URI: ",config.ControllerBrokerAddress )

	desiredTopic  := "#"
	subscriberConfig := models.NewMqttClientConfiguration(config.ControllerBrokerAddress,"domainController")
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
	config := configuration.NewDomainControllerConfiguration(controllerID, parentDomain,ownDomain, brokerURL,gatewayBrokerURL)
	configuration.NewDomainControllerConfigurationManager().StoreConfig(config)
}
