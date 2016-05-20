package main

import (

	"github.com/tkrex/IDS/daemon/layers"
	"time"
	"github.com/tkrex/IDS/common/models"
	"github.com/tkrex/IDS/common/layers"
)

func main() {

	go startBrokerRegistration()
	go startControlMessageProcessing()
	go startTopicProcessing()
	for {
		time.Sleep(time.Second)
	}
}


func startTopicProcessing() {
	brokerAddress := "tcp://localhost:1883"
	desiredTopic  := "#"
	subscriberConfig := models.NewMqttClientConfiguration(brokerAddress,desiredTopic,"subscriber")
	subscriber := common.NewMqttSubscriber(subscriberConfig,true)
	topicProcessor := layers.NewTopicProcessor(subscriber.IncomingTopicsChannel())


	_ = layers.NewDomainInformationForwarder(topicProcessor.ForwardSignalChannel())

}
func startControlMessageProcessing() {
	brokerAddress := "tcp://localhost:1883"
	desiredTopic  := "ControlMessage"
	//TODO: figure out client id
	subscriberConfig := models.NewMqttClientConfiguration(brokerAddress,desiredTopic,"controlMessageSubscriber")
	subscriber := common.NewMqttSubscriber(subscriberConfig,true)
	_ = common.NewControlMessageProcessor(subscriber.IncomingTopicsChannel())

}

func startBrokerRegistration() {
	_ = layers.NewBrokerRegistrationWorker()
}
