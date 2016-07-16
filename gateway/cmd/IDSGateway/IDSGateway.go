package main

import (
	"github.com/tkrex/IDS/gateway/providing"
	"github.com/tkrex/IDS/gateway/domainControllerManagement"
	"github.com/tkrex/IDS/common/models"
	"fmt"
	"github.com/tkrex/IDS/gateway/scaling"
)

func main() {


	_ = domainControllerManagement.NewDomainContollerManagementInterface("8080")

	providing.NewIDSGatewayWebInterface("8000")

	scalingManager := scaling.NewDockerManager()
	if domainController,_ := scalingManager.StartDomainControllerInstance(models.NewRealWorldDomain("none"),models.NewRealWorldDomain("default")); domainController != nil {
		fmt.Println("Default DomainController created: ", domainController)
	}
	for  {}
}

