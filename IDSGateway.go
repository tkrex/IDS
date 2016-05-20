package main

import (
	"github.com/tkrex/IDS/gateway/providing"
	"github.com/tkrex/IDS/gateway/controlling"
)

func main() {

	providing.NewIDSGatewayWebInterface("8080")
	webInterface := controlling.NewServerMaintenanceWebInterface("8080")
	_ = controlling.NewControlMessageProcessor(webInterface.IncomingControlMessagesChannel())
	for  {}
}

