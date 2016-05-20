package main

import (
	"github.com/tkrex/IDS/gateway/providing"
)

func main() {

	providing.NewIDSGatewayWebInterface("8080")
	//webInterface := controlling.NewServerMaintenanceWebInterface("8080")
	//_ = controlling.NewControlMessageProcessor(webInterface.IncomingControlMessagesChannel())
	for  {}
}

