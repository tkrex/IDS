package main

import (
	"time"
	"github.com/tkrex/IDS/common/models"
	"github.com/tkrex/IDS/common/subscribing"
	"github.com/tkrex/IDS/daemon/processing"
	"github.com/tkrex/IDS/daemon/forwarding"
	"github.com/tkrex/IDS/daemon/registration"
	"github.com/tkrex/IDS/daemon/configuration"
	"fmt"
)

func main() {
	config := configuration.DaemonConfigurationManagerInstance().Config()
	fmt.Println("Daemon Config ",config)
	//startBrokerRegistration(config)
	startTopicProcessing(config)
	for {
		time.Sleep(time.Second)
	}
}

func startBrokerRegistration(config *configuration.DaemonConfiguration) {
	_ = registration.NewBrokerRegistrattionManager(config.RegistrationURL)
}

func startTopicProcessing(config *configuration.DaemonConfiguration) {
	desiredTopic := "#"
	subscriberConfig := models.NewMqttClientConfiguration(config.BrokerURL, "subscriber")
	subscriber := subscribing.NewMqttSubscriber(subscriberConfig, desiredTopic, false)
	topicProcessor := processing.NewTopicProcessor(subscriber.IncomingTopicsChannel())
	_ = forwarding.NewDomainInformationForwarder(topicProcessor.ForwardSignalChannel())
}


