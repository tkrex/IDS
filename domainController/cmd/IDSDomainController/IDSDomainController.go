package main

import (

	"github.com/tkrex/IDS/common/models"
	"github.com/tkrex/IDS/common/subscribing"
	"github.com/tkrex/IDS/domainController/processing"
	"github.com/tkrex/IDS/domainController/providing"
	"os"
	"fmt"
	"github.com/tkrex/IDS/common/controlling"
)



func main() {

	go startDomainInformationProcessing()
	go startControlMessageProcessing()
	_ = providing.NewDomainInformationRESTProvider("8080")
	for {}
}


func startDomainInformationProcessing() {
	//producer layer

	brokerAddress := os.Getenv("BROKER_URI")
	fmt.Println("Broker URI: ",brokerAddress)
	port := "1883"
	var protocol models.MqttProtocol  = "tcp"
	desiredTopic  := "DomainInformation"
	subscriberConfig := models.NewMqttClientConfiguration(brokerAddress,port,protocol,desiredTopic,"domainController")
	subscriber := subscribing.NewMqttSubscriber(subscriberConfig,false)
	//processing layer
	_ = processing.NewDomainInformationProcessor(subscriber.IncomingTopicsChannel())
}

func startControlMessageProcessing() {
	brokerAddress := os.Getenv("MANAGEMENT_BROKER_URL")
	desiredTopic  := "ControlMessage"
	//TODO: figure out client id
	subscriberConfig := models.NewMqttClientConfiguration(brokerAddress,"1883","tcp",desiredTopic,"controlMessageSubscriber")
	subscriber := subscribing.NewMqttSubscriber(subscriberConfig,true)
	_ = controlling.NewControlMessageProcessor(subscriber.IncomingTopicsChannel())

}
