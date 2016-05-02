package main

import (

	"time"
	"github.com/tkrex/IDS/common/layers"
	"github.com/tkrex/IDS/common/models"
	"github.com/tkrex/IDS/domainController"
)



func main() {


	//processing layer
	topicChannel := make(chan *models.RawTopicMessage,100)
	_ = domainController.NewDomainInformationProcessor(topicChannel)
	//producer layer
	brokerAddress := "tcp://localhost:1883"
	desiredTopic  := "domainController"
	subscriberConfig := models.NewMqttClientConfiguration(brokerAddress,desiredTopic,"domainController")
	subscriber := common.NewMqttSubscriber(subscriberConfig,topicChannel,false)

	_ = domainController.NewRestInformationProvider("8080")
	for {
		time.Sleep(time.Second)
	}
	subscriber.Close()


	//worker := common.NewWebsiteCategorizationWorker("owLf4fHmY0jMwQLNapZD","http://api.webshrinker.com/categories/v2")
	//worker.RequestCategoriesForWebsite("www1.in.tum.de")

}

