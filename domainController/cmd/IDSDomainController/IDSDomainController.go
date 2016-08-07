package main

import (

	"github.com/tkrex/IDS/common/models"
	"github.com/tkrex/IDS/common/subscribing"
	"github.com/tkrex/IDS/domainController/processing"
	"github.com/tkrex/IDS/domainController/providing"
	"fmt"
	"github.com/tkrex/IDS/domainController/configuration"
	"github.com/tkrex/IDS/domainController/forwarding"
)



func main() {
	controllerConfig  := configuration.DomainControllerConfigurationManagerInstance().Config()

	isTopLevelController := controllerConfig.ParentDomain.Name == "none"
	isDefaultController := controllerConfig.OwnDomain.Name == "default"
	shouldForwardDomainInformation := isDefaultController || !isTopLevelController
	go startDomainInformationProcessing(shouldForwardDomainInformation)
	if isTopLevelController {
		_ = providing.NewDomainInformationRESTProvider("8080")
	}
	for {}
}


func startDomainInformationProcessing(forwardFlag bool) {
	//producer layer


	config := configuration.DomainControllerConfigurationManagerInstance().Config()

	fmt.Println("Config: ",config)

	desiredTopic  := "#"
	subscriberConfig := models.NewMqttClientConfiguration(config.ControllerBrokerAddress,config.DomainControllerID)
	subscriber := subscribing.NewMqttSubscriber(subscriberConfig,desiredTopic)
	//processing layer
	processor := processing.NewDomainInformationProcessor(subscriber.IncomingTopicsChannel(),forwardFlag)
	if forwardFlag {
		forwarding.NewDomainInformationForwarder(processor.ForwardSignalChannel(),processor.RoutingInformationChannel())
	}

}

