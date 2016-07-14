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
	controllerConfig  := configuration.DomainControllerConfigurationManager().Config()

	isTopLevelController := controllerConfig.ParentDomain == nil
	go startDomainInformationProcessing(!isTopLevelController)
	if isTopLevelController {
		_ = providing.NewDomainInformationRESTProvider("8080")
	}
	forwarding.NewDomainInformationForwarder()
	for {}
}


func startDomainInformationProcessing(forwardFlag bool) {
	//producer layer


	config := configuration.DomainControllerConfigurationManager().Config()

	fmt.Println("Config: ",config)

	desiredTopic  := "#"
	subscriberConfig := models.NewMqttClientConfiguration(config.ControllerBrokerAddress,config.DomainControllerID)
	subscriber := subscribing.NewMqttSubscriber(subscriberConfig,desiredTopic,false)
	//processing layer
	processor := processing.NewDomainInformationProcessor(subscriber.IncomingTopicsChannel(),forwardFlag)
	forwarding.NewDomainInformationForwarder(processor.ForwardSignalChannel())
	
}

