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
)

func main() {

	//go startBrokerRegistration()
	//go startControlMessageProcessing()
	go startTopicProcessing()
	startControlInterface()
	for {
		time.Sleep(time.Second)
	}
}


func startTopicProcessing() {
	brokerAddress := os.Getenv("BROKER_URI")
	fmt.Println("Broker URL", brokerAddress)
	desiredTopic  := "#"
	subscriberConfig := models.NewMqttClientConfiguration(brokerAddress,"1883","tcp",desiredTopic,"subscriber")
	subscriber := subscribing.NewMqttSubscriber(subscriberConfig,true)
	topicProcessor := processing.NewTopicProcessor(subscriber.IncomingTopicsChannel())


	_ = forwarding.NewDomainInformationForwarder(topicProcessor.ForwardSignalChannel())

}
func startControlMessageProcessing() {
	brokerAddress := os.Getenv("MANAGEMENT_BROKER_URL")
	desiredTopic  := "ControlMessage"
	//TODO: figure out client id
	subscriberConfig := models.NewMqttClientConfiguration(brokerAddress,"1883","tcp",desiredTopic,"controlMessageSubscriber")
	subscriber := subscribing.NewMqttSubscriber(subscriberConfig,true)
	_ = controlling.NewControlMessageProcessor(subscriber.IncomingTopicsChannel())
}

func startBrokerRegistration() {
	_ = registration.NewBrokerRegistrationWorker(os.Getenv("MANAGEMENT_INTERFACE_URL"))
}

func startControlInterface () {
	_ = configuration.NewConfigurationInterface("8080")
}
