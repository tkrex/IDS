package main

import (

)
import (
	"time"
	"flag"
	"github.com/tkrex/IDS/common/layers"
)

func main() {
	flag.Parse()
	brokerAddress := "tcp://localhost:1883"
	desiredTopic  := "#"
	var subscriber common.InformationProducer
	subscriber = common.NewMqttSubscriber(brokerAddress,desiredTopic)

	time.Sleep(time.Second * 60)
	subscriber.Close()
}
