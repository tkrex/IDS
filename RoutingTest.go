package main

import ()
import (
	"github.com/tkrex/IDS/common/models"
	"fmt"
)

func main() {
	//domainControllerManager := controlling.NewDomainControllerManager()
	testDomain := models.NewRealWorldDomain("Test/1/2/3")
	domainLevels := testDomain.DomainLevels()
	for i := len(domainLevels) - 1; i >= 0; i-- {
		fmt.Println("Searching Domain Controller for domain: ", testDomain)
		//requestedDomainController= dbWorker.FindDomainControllerForDomain(domain)
		testDomain = testDomain.ParentDomain()
	}

}
