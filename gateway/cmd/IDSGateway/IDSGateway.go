package main

import (
	"github.com/tkrex/IDS/gateway/domainControllerManagement"
	"net/url"
)

func main() {
	managementBrokerAddress,_ := url.Parse("ws://localhost:11883")


	_ = domainControllerManagement.NewServerMaintenanceWebInterface("8080",managementBrokerAddress)
	//managementServerAddress := "http://localhost:8080"

	//providing.NewIDSGatewayWebInterface("8000")



	//worker := common.NewDomainControllerRegistrationWorker(managementServerAddress)
	//worker.RequestNewDomainControllerForDomain(models.NewRealWorldDomain("default"))

	for  {}
}

