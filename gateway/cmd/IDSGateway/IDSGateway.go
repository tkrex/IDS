package main

import (
	"github.com/tkrex/IDS/gateway/providing"
	"github.com/tkrex/IDS/gateway/domainControllerManagement"

)

func main() {


	_ = domainControllerManagement.NewDomainContollerManagementInterface("8080")

	providing.NewIDSGatewayWebInterface("8000")

/*	domainControllerManager := domainControllerManagement.NewDomainControllerManager()
	if domainController,_ := domainControllerManager.StartNewDomainControllerInstance(models.NewRealWorldDomain("default"),models.NewRealWorldDomain("none")); domainController != nil {
		fmt.Println("Default DomainController created: ", domainController)
	}*/
	for  {}
}

