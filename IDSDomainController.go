package main

import (

	"github.com/tkrex/IDS/common/models"
	"github.com/tkrex/IDS/common/subscribing"
	"github.com/tkrex/IDS/domainController/processing"
)



func main() {

	startDomainInformationProcessing()
	for {}
}


func startDomainInformationProcessing() {
	//producer layer
	brokerAddress := "tcp://localhost:1883"
	desiredTopic  := "DomainInformation"
	subscriberConfig := models.NewMqttClientConfiguration(brokerAddress,desiredTopic,"domainController")
	subscriber := subscribing.NewMqttSubscriber(subscriberConfig,false)
	//processing layer
	_ = processing.NewDomainInformationProcessor(subscriber.IncomingTopicsChannel())
}

