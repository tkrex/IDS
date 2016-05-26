package main

import (
	"github.com/tkrex/IDS/common/models"
	"github.com/tkrex/IDS/common/publishing"
)

func main() {
	publishConfig := models.NewMqttClientConfiguration("localhost","1883","ws","test","testClient")
	publisher := publishing.NewMqttPublisher(publishConfig)
	publisher.Publish([]byte("Test"))
}
