package main

import (

)
import (
	"github.com/tkrex/IDS/gateway"
)

func main() {

	gateway.NewIDSGatewayWebInterface("8080")
	webInterface := gateway.NewServerMaintenanceWebInterface("8080")
	_ = gateway.NewControlMessageProcessor(webInterface.IncomingControlMessagesChannel())
	for  {}
}

