package main

import (
	"github.com/tkrex/IDS/gateway/clusterManagement"

	"github.com/tkrex/IDS/common/models"
	"fmt"
	"github.com/tkrex/IDS/gateway/InformationRequestManagement"
	"github.com/tkrex/IDS/gateway/brokerRegistration"
)

//Starts up the component of the Gateway
func main() {

	domainControllerManager := clusterManagement.NewClusterManager()

	if domainController,_ := domainControllerManager.StartNewDomainControllerInstance(models.NewRealWorldDomain("default"),models.NewRealWorldDomain("none")); domainController != nil {
		fmt.Println("Default DomainController created: ", domainController)
	}
	clusterManagement.NewClusterManagerInterface("8080")
	brokerRegistration.NewBrokerRegistrationInterface("8001")

	informationRequestManagement.NewInformationRequestInterface("8000")
	for  {}
}

