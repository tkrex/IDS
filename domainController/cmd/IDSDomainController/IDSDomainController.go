package main

import (

	"github.com/tkrex/IDS/common/models"
	"github.com/tkrex/IDS/common/subscribing"
	"github.com/tkrex/IDS/domainController/processing"
	"github.com/tkrex/IDS/domainController/providing"
	"os"
	"fmt"
	"github.com/tkrex/IDS/common/controlling"
	"net/url"
)



func main() {

	go startDomainInformationProcessing()
	go startControlMessageProcessing()
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

func startControlMessageProcessing() {
	brokerAddressString := os.Getenv("MANAGEMENT_BROKER_URL")
	if brokerAddressString == "" {
		brokerAddressString = "ws://localhost:11883"
	}
	brokerAddress, _ := url.Parse(brokerAddressString)

	desiredTopic  := "ControlMessage"
	//TODO: figure out client id
	subscriberConfig := models.NewMqttClientConfiguration(brokerAddress,"controlMessageSubscriber")
	subscriber := subscribing.NewMqttSubscriber(subscriberConfig,desiredTopic,true)
	_ = controlling.NewControlMessageProcessor(subscriber.IncomingTopicsChannel())

}
