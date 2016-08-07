package main

import (
	"github.com/tkrex/IDS/common/models"
	"time"
	"encoding/json"
	"net/url"
	"github.com/tkrex/IDS/common/publishing"
)

func main() {



	broker := models.NewBroker()
	topics := []*models.TopicInformation{}

	broker = models.NewBroker()
	broker.ID = "weatherBroker"
	broker.IP = "12.12.12.12:1833"
	broker.InternetDomain = "krex.in.tum.de"
	broker.Statistics.NumberOfTopics = 1022
	broker.Statistics.ReceivedTopicsPerSeconds = 10
	broker.RealWorldDomain = models.NewRealWorldDomain("weather")
	broker.Geolocation = models.NewGeolocation("Germany", "Bavaria", "Garching", 11.6309, 48.2499)

	topic := models.NewTopicInformation("/fmi/server-room", "{\"temperature\":-6}", time.Now())
	topic.UpdateBehavior.AverageUpdateIntervalInSeconds = 180
	topic.UpdateBehavior.UpdateIntervalDeviation = 3.0
	topic.PayloadSimilarity = 80.5
	topic.UpdateBehavior.Reliability = "automatic"
	topics = append(topics, topic)

	topic = models.NewTopicInformation("/fmi/ls1", "{\"temperature\":30}", time.Now())
	topic.UpdateBehavior.AverageUpdateIntervalInSeconds = 120
	topic.UpdateBehavior.UpdateIntervalDeviation = 70.0
	topic.PayloadSimilarity = 87
	topic.UpdateBehavior.Reliability = "semi-automatic"
	topics = append(topics, topic)

	topic = models.NewTopicInformation("/fmi/ls2", "{\"temperature\":10}", time.Now())
	topic.UpdateBehavior.AverageUpdateIntervalInSeconds = 30
	topic.UpdateBehavior.UpdateIntervalDeviation = 200
	topic.PayloadSimilarity = 40
	topic.UpdateBehavior.Reliability = "non-deterministic"
	topics = append(topics, topic)

	topic = models.NewTopicInformation("/fmi/smart-lab", "{\"temperature\":10}", time.Now())
	topic.UpdateBehavior.AverageUpdateIntervalInSeconds = 60
	topic.UpdateBehavior.UpdateIntervalDeviation = 5
	topic.PayloadSimilarity = 99
	topic.UpdateBehavior.Reliability = "automatic"
	topics = append(topics, topic)

	message := models.NewDomainInformationMessage(broker.RealWorldDomain, broker, topics)

	bytes,_ := json.Marshal(message)

	brokerURL,_  :=  url.Parse("ws://10.40.53.21:32870")
	publishConfig := models.NewMqttClientConfiguration(brokerURL,"testPublisher")
	publisher := publishing.NewMqttPublisher(publishConfig,false)
	publisher.Publish(bytes,"weatherBroker")
	publisher.Close()
}
