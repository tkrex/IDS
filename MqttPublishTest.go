package main

import (
	"github.com/tkrex/IDS/common/models"
	"github.com/tkrex/IDS/common/publishing"
)

func main() {
	publishConfig := models.NewMqttClientConfiguration("localhost","11883","ws","test","testClient")
	publisher := publishing.NewMqttPublisher(publishConfig,false)
	publisher.Publish([]byte("Test"))
}
