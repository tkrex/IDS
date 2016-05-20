package main

import (

	"time"
	"github.com/tkrex/IDS/common/layers"
	"github.com/tkrex/IDS/common/models"
	"github.com/tkrex/IDS/domainController"
)



func main() {



	startDomainInformationProcessing()
	startControlMessageProcessing()
	for {}
}


func startDomainInformationProcessing() {
	//producer layer
	brokerAddress := "tcp://localhost:1883"
	desiredTopic  := "DomainInformation"
	subscriberConfig := models.NewMqttClientConfiguration(brokerAddress,desiredTopic,"domainController")
	subscriber := common.NewMqttSubscriber(subscriberConfig,false)
	//processing layer
	_ = domainController.NewDomainInformationProcessor(subscriber.IncomingTopicsChannel())
}

