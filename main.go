package main

import (

)
import (
	"time"
	"flag"
	"github.com/tkrex/IDS/common/layers"
	"github.com/tkrex/IDS/daemon/layers"
)

func main() {
	flag.Parse()
	brokerAddress := "tcp://localhost:1883"
	desiredTopic  := "topicInformation"
	var subscriber common.InformationProducer
	subscriber = common.NewMqttSubscriber(brokerAddress,desiredTopic)
	publisher := layers.NewMqttPublisher(brokerAddress, desiredTopic)

	time.Sleep(time.Second * 60)
	subscriber.Close()
	publisher.Close()
}
