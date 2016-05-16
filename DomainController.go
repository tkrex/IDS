package main

import (

	"time"
	"github.com/tkrex/IDS/common/layers"
	"github.com/tkrex/IDS/common/models"
	"github.com/tkrex/IDS/domainController"
)



func main() {



	//producer layer
	brokerAddress := "tcp://localhost:1883"
	desiredTopic  := "domainController"
	subscriberConfig := models.NewMqttClientConfiguration(brokerAddress,desiredTopic,"domainController")
	subscriber := common.NewMqttSubscriber(subscriberConfig,false)
	//processing layer
	_ = domainController.NewDomainInformationProcessor(subscriber.IncomingTopicsChannel())

	for {
		time.Sleep(time.Second)
	}
	subscriber.Close()


}

