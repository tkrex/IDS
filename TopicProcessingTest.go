package main

import (
	"github.com/tkrex/IDS/common/models"
	"github.com/tkrex/IDS/common/publishing"
	"net/url"
	"encoding/json"
	"time"
	"fmt"
)

func main() {
	brokerAddress,_ := url.Parse("ws://localhost:11883")
	publishConfig := models.NewMqttClientConfiguration(brokerAddress,"testClient")
	publisher := publishing.NewMqttPublisher(publishConfig,false)


	for i := 0;i<100;i++ {
		fmt.Println("Message: " , i)
		payload := make(map[string]int)
		payload["id"] = i
		payload["temp"] = 30
		json,_ := json.Marshal(&payload)
		publisher.Publish(json,"testTopic")
		time.Sleep(1 * time.Second)
	}
	publisher.Close()
}
