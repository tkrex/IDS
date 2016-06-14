package main

import (
	"time"
	"github.com/tkrex/IDS/common/models"
	"github.com/tkrex/IDS/common/subscribing"
	"github.com/tkrex/IDS/daemon/processing"
	"github.com/tkrex/IDS/daemon/forwarding"
	"github.com/tkrex/IDS/common/controlling"
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
	startControlInterface()
	for {
		time.Sleep(time.Second)
	}
}


func startTopicProcessing() {
	brokerAddressString := os.Getenv("BROKER_URI")
	brokerAddress, _ := url.Parse(brokerAddressString)
	fmt.Println("Broker URL", brokerAddress)
	desiredTopic  := "#"
	subscriberConfig := models.NewMqttClientConfiguration(brokerAddress,"subscriber")
	subscriber := subscribing.NewMqttSubscriber(subscriberConfig,desiredTopic,false)
	topicProcessor := processing.NewTopicProcessor(subscriber.IncomingTopicsChannel())


	_ = forwarding.NewDomainInformationForwarder(topicProcessor.ForwardSignalChannel())

}
func startControlMessageProcessing() {
	brokerAddressString := os.Getenv("MANAGEMENT_BROKER_URL")
	brokerAddress, _ := url.Parse(brokerAddressString)

	desiredTopic  := "ControlMessage"
	//TODO: figure out client id
	subscriberConfig := models.NewMqttClientConfiguration(brokerAddress,"controlMessageSubscriber")
	subscriber := subscribing.NewMqttSubscriber(subscriberConfig,desiredTopic,false)
	_ = controlling.NewControlMessageProcessor(subscriber.IncomingTopicsChannel())
}

func startBrokerRegistration() {
	_ = registration.NewBrokerRegistrationWorker(os.Getenv("MANAGEMENT_INTERFACE_URL"))
}

func startControlInterface () {
	_ = configuration.NewConfigurationInterface("8080")
}
