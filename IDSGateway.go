package main

import (
	"github.com/tkrex/IDS/gateway/providing"
	"github.com/tkrex/IDS/gateway/controlling"
	"github.com/tkrex/IDS/common"
	"github.com/tkrex/IDS/common/models"
)

func main() {

	providing.NewIDSGatewayWebInterface("8080")

	_ = controlling.NewServerMaintenanceWebInterface("8000")
	worker := common.NewDomainControllerRegistrationWorker()
	worker.RequestNewDomainControllerForDomain(models.NewRealWorldDomain("default"))

	//_ = controlling.NewControlMessageProcessor(webInterface.IncomingControlMessagesChannel())
	for  {}
}

