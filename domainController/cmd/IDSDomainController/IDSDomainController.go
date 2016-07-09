package main

import (

	"github.com/tkrex/IDS/common/models"
	"github.com/tkrex/IDS/common/subscribing"
	"github.com/tkrex/IDS/domainController/processing"
	"github.com/tkrex/IDS/domainController/providing"
	"fmt"
	"github.com/tkrex/IDS/domainController/configuration"
)



func main() {
	controllerConfig  := configuration.NewDomainControllerConfigurationManager().InitConfig()

	isTopLevelController := controllerConfig.ParentDomain == nil
	go startDomainInformationProcessing(!isTopLevelController)
	if isTopLevelController {
		_ = providing.NewDomainInformationRESTProvider("8080")
	}
	for {}
}


func startDomainInformationProcessing(forwardFlag bool) {
	//producer layer


	config, _ := configuration.NewDomainControllerConfigurationManager().DomainControllerConfig()

	fmt.Println("Broker URI: ",config.ControllerBrokerAddress )

	desiredTopic  := "#"
	subscriberConfig := models.NewMqttClientConfiguration(config.ControllerBrokerAddress,"domainController")
	subscriber := subscribing.NewMqttSubscriber(subscriberConfig,desiredTopic,false)
	//processing layer
	_ = processing.NewDomainInformationProcessor(subscriber.IncomingTopicsChannel(),forwardFlag)
}

