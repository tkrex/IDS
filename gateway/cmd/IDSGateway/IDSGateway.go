package main

import (
	"github.com/tkrex/IDS/gateway/providing"
	"github.com/tkrex/IDS/gateway/domainControllerManagement"

	"github.com/tkrex/IDS/common/models"
	"fmt"
)

func main() {

	domainControllerManager := domainControllerManagement.NewDomainControllerManager()
	if domainController,_ := domainControllerManager.StartNewDomainControllerInstance(models.NewRealWorldDomain("default"),models.NewRealWorldDomain("none")); domainController != nil {
		fmt.Println("Default DomainController created: ", domainController)
	}
	_ = domainControllerManagement.NewDomainContollerManagementInterface("8080")

	providing.NewIDSGatewayWebInterface("8000")
	for  {}
}

