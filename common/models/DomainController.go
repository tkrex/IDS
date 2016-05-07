package models

import "fmt"

type DomainController  struct{
	IpAddress string `json:"ipAddress" bson:"ipAddress"`
	Domain    *RealWorldDomain `json:"domain" bson:"domain"`
}

func NewDomainController(ipAddress string, domain *RealWorldDomain) *DomainController {
	controller := new(DomainController)
	controller.IpAddress = ipAddress
	controller.Domain = domain
	return controller
}

func (controller *DomainController) String() string {
	return fmt.Sprintf("IP-Address: %s, Domain: %s",controller.IpAddress, controller.Domain)
}