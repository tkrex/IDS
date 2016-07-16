package main

import (
	"github.com/tkrex/IDS/gateway/providing"
	"github.com/tkrex/IDS/gateway/domainControllerManagement"
	"github.com/tkrex/IDS/common"
	"github.com/tkrex/IDS/common/models"
	"fmt"
)

func main() {


	_ = domainControllerManagement.NewDomainContollerManagementInterface("8080")

	providing.NewIDSGatewayWebInterface("8000")

	scalingManager := common.NewScalingManager()
	if domainController := scalingManager.CreateNewDominControllerForDomain(models.NewRealWorldDomain("default")); domainController != nil {
		fmt.Println("Default DomainController created: ", domainController)
	}
	for  {}
}

