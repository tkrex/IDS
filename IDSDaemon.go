package main

import (

	"time"
	"github.com/tkrex/IDS/common/layers"
	"github.com/tkrex/IDS/common/models"
)

func main() {


	// persistence layer
	persistenceMananger := common.NewMemoryPersistenceManager()
	//processing layer
	topicChannel := make(chan *models.Topic,100)
	_ = common.NewTopicProcessor(persistenceMananger,topicChannel)
	//producer layer
	brokerAddress := "tcp://localhost:1883"
	desiredTopic  := "#"
	subscriberConfig := models.NewMqttClientConfiguration(brokerAddress,desiredTopic,"subscriber")
	subscriber := common.NewMqttSubscriber(subscriberConfig,topicChannel)

	_ = common.NewTopicForwarder(persistenceMananger,time.Second *10)

	time.Sleep(time.Second * 60)
	subscriber.Close()


	//worker := common.NewWebsiteCategorizationWorker("owLf4fHmY0jMwQLNapZD","http://api.webshrinker.com/categories/v2")
	//worker.RequestCategoriesForWebsite("www1.in.tum.de")

}

