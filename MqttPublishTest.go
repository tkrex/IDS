package main

import (
	"github.com/tkrex/IDS/common/models"
	"github.com/tkrex/IDS/common/publishing"
	"time"
	"encoding/json"
)

func main() {


	//DEBUG CODE
	domain := models.NewRealWorldDomain("education")
	broker := models.NewBroker()
	broker.ID = "testID"
	broker.IP = "new host"
	broker.RealWorldDomain = models.NewRealWorldDomain("education")
	broker.Geolocation = models.NewGeolocation("germany","bavaria","munich",11.6309,48.2499)

	topics := []*models.Topic{}

	for i := 0; i < 5; i++ {
		topic := models.NewTopic("/home/kitchen","{\"temperature\":3}",time.Now())
		topic.UpdateBehavior.UpdateIntervalDeviation = 3.0
		topics = append(topics, topic)
	}


	message := models.NewDomainInformationMessage(domain,broker,topics)

	json,_ := json.Marshal(message)

	publishConfig := models.NewMqttClientConfiguration("localhost","11883","ws","testID","testClient")
	publisher := publishing.NewMqttPublisher(publishConfig,false)
	publisher.Publish(json)
}
