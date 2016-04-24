package main

import (

	"time"
	"github.com/tkrex/IDS/common/layers"
)

func main() {
	brokerAddress := "tcp://localhost:1883"
	desiredTopic  := "#"
	var subscriber common.InformationProducer
	subscriber = common.NewMqttSubscriber(brokerAddress,desiredTopic)

	time.Sleep(time.Second * 60)
	subscriber.Close()


	//worker := common.NewWebsiteCategorizationWorker("owLf4fHmY0jMwQLNapZD","http://api.webshrinker.com/categories/v2")
	//worker.RequestCategoriesForWebsite("www1.in.tum.de")

}

