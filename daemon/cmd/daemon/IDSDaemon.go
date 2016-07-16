package main

import (
	"time"
	"github.com/tkrex/IDS/common/models"
	"github.com/tkrex/IDS/common/subscribing"
	"github.com/tkrex/IDS/daemon/processing"
	"github.com/tkrex/IDS/daemon/forwarding"
	"github.com/tkrex/IDS/daemon/registration"
	"github.com/tkrex/IDS/daemon/configuration"
	"os"
	"fmt"
	"net/url"
)

func main() {

	startBrokerRegistration()

	//go startControlMessageProcessing()
	go startTopicProcessing()
	//startControlInterface()
	for {
		time.Sleep(time.Second)
	}
}

func startTopicProcessing() {
	brokerAddressString := os.Getenv("BROKER_URI")
	if brokerAddressString == "" {
		brokerAddressString = "ws://localhost:11883"
	}
	brokerAddress, _ := url.Parse(brokerAddressString)
	fmt.Println("Broker URL", brokerAddress)

	desiredTopic := "#"
	subscriberConfig := models.NewMqttClientConfiguration(brokerAddress, "subscriber")
	subscriber := subscribing.NewMqttSubscriber(subscriberConfig, desiredTopic, false)
	topicProcessor := processing.NewTopicProcessor(subscriber.IncomingTopicsChannel())

	_ = forwarding.NewDomainInformationForwarder(topicProcessor.ForwardSignalChannel())
}

func startBrokerRegistration() {
	registrationServiceUrlString := os.Getenv("MANAGEMENT_INTERFACE_URL")
	if registrationServiceUrlString == "" {
		registrationServiceUrlString = "http://localhost:8000"
	}
	registrationServiceUrl,_ := url.Parse(registrationServiceUrlString)
	_ = registration.NewBrokerRegistrationWorker(registrationServiceUrl)
}

func startControlInterface () {
	_ = configuration.NewConfigurationInterface("8080")
}
