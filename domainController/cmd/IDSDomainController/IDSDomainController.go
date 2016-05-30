package main

import (

	"github.com/tkrex/IDS/common/models"
	"github.com/tkrex/IDS/common/subscribing"
	"github.com/tkrex/IDS/domainController/processing"
	"github.com/tkrex/IDS/domainController/providing"
)



func main() {

	startDomainInformationProcessing()
	_ = providing.NewDomainInformationRESTProvider("8080")
	for {}
}


func startDomainInformationProcessing() {
	//producer layer
	brokerAddress := "localhost"
	port := "1883"
	var protocol models.MqttProtocol  = "tcp"
	desiredTopic  := "DomainInformation"
	subscriberConfig := models.NewMqttClientConfiguration(brokerAddress,port,protocol,desiredTopic,"domainController")
	subscriber := subscribing.NewMqttSubscriber(subscriberConfig,false)
	//processing layer
	_ = processing.NewDomainInformationProcessor(subscriber.IncomingTopicsChannel())
}

