package main

import (
	"github.com/tkrex/IDS/gateway/providing"
	"github.com/tkrex/IDS/gateway/controlling"
	"github.com/tkrex/IDS/common"
	"github.com/tkrex/IDS/common/models"
)

func main() {


	managementServerAddress := "http://localhost"

	providing.NewIDSGatewayWebInterface("8000")
	managementBrokerAddress := "localhost"

	_ = controlling.NewServerMaintenanceWebInterface("8080",managementBrokerAddress)

	worker := common.NewDomainControllerRegistrationWorker(managementServerAddress)
	worker.RequestNewDomainControllerForDomain(models.NewRealWorldDomain("default"))

	for  {}
}

