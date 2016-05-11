package main

import (

	"github.com/tkrex/IDS/daemon/layers"
	"time"
	"github.com/tkrex/IDS/common/models"
	"github.com/tkrex/IDS/common/layers"
)

func main() {

	// persistence layer
	//processing layer
	//producer layer
	brokerAddress := "tcp://localhost:1883"
	desiredTopic  := "#"
	subscriberConfig := models.NewMqttClientConfiguration(brokerAddress,desiredTopic,"subscriber")
	subscriber := common.NewMqttSubscriber(subscriberConfig,true)
	topicProcessor := layers.NewTopicProcessor(subscriber.IncomingTopicsChannel())


	_ = layers.NewDomainInformationForwarder(topicProcessor.ForwardSignalChannel())

	for {
		time.Sleep(time.Second)
	}
	subscriber.Close()


	//worker := common.NewWebsiteCategorizationWorker("owLf4fHmY0jMwQLNapZD","http://api.webshrinker.com/categories/v2")
	//worker.RequestCategoriesForWebsite("www1.in.tum.de")


}

