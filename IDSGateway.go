package main

import (

)
import (
	"github.com/tkrex/IDS/gateway"
)

func main() {


//_ = gateway.NewControlMessageForwarder()
	_ = gateway.NewIDSGatewayWebInterface("8080")
}

