package main

import (

	"time"
	"github.com/tkrex/IDS/common/layers"
	"github.com/tkrex/IDS/common/models"
	"github.com/tkrex/IDS/domainController"
)



func main() {

	// persistence layer
	var persistenceManager := domainController.NewDomainInformationMemoryPersistenceManager()
	//processing layer
	topicChannel := make(chan *models.RawTopicMessage,100)
	_ = domainController.NewDomainInformationProcessor(persistenceManager,topicChannel)
	//producer layer
	brokerAddress := "tcp://localhost:1883"
	desiredTopic  := "domainController"
	subscriberConfig := models.NewMqttClientConfiguration(brokerAddress,desiredTopic,"domainController")
	subscriber := common.NewMqttSubscriber(subscriberConfig,topicChannel,false)

	_ = domainController.NewRestInformationProvider(persistenceManager,"8080")
	for {
		time.Sleep(time.Second)
	}
	subscriber.Close()


	//worker := common.NewWebsiteCategorizationWorker("owLf4fHmY0jMwQLNapZD","http://api.webshrinker.com/categories/v2")
	//worker.RequestCategoriesForWebsite("www1.in.tum.de")

}

