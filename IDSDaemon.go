package main

import (

	"github.com/tkrex/IDS/common/models"
	"github.com/tkrex/IDS/common/layers"
	"github.com/tkrex/IDS/daemon/layers"
	"time"
)

func main() {

	// persistence layer
	//processing layer
	topicChannel := make(chan *models.RawTopicMessage,100)
	_ = layers.NewTopicProcessor(topicChannel)
	//producer layer
	brokerAddress := "tcp://localhost:1883"
	desiredTopic  := "#"
	subscriberConfig := models.NewMqttClientConfiguration(brokerAddress,desiredTopic,"subscriber")
	subscriber := common.NewMqttSubscriber(subscriberConfig,topicChannel,true)

	_ = layers.NewTopicForwarder(time.Second *10)

	for {
		time.Sleep(time.Second)
	}
	subscriber.Close()


	//worker := common.NewWebsiteCategorizationWorker("owLf4fHmY0jMwQLNapZD","http://api.webshrinker.com/categories/v2")
	//worker.RequestCategoriesForWebsite("www1.in.tum.de")


}

