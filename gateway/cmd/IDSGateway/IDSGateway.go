package main

import (
	"github.com/tkrex/IDS/gateway/clusterManagement"

	"github.com/tkrex/IDS/common/models"
	"fmt"
	"github.com/tkrex/IDS/gateway/brokerRegistration"
	"github.com/tkrex/IDS/gateway/InformationRequestManagement"
)

func main() {

	domainControllerManager := clusterManagement.NewClusterManager()

	if domainController,_ := domainControllerManager.StartNewDomainControllerInstance(models.NewRealWorldDomain("default"),models.NewRealWorldDomain("none")); domainController != nil {
		fmt.Println("Default DomainController created: ", domainController)
	}
	clusterManagement.NewClusterManagerInterface("8080")
	registration.NewBrokerRegistrationInterface("8001")

	informationRequestManagement.InformationRequestInterface("8000")
	for  {}
}

