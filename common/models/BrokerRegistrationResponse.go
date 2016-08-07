package models

import (
	"fmt"
)

type BrokerRegistrationResponse struct {
	Broker           *Broker `json:"broker"`
	DomainController *DomainController `json:"domainControllers"`
}

func NewBrokerRegistrationResponse(broker *Broker, controller *DomainController) *BrokerRegistrationResponse {
	response := new(BrokerRegistrationResponse)
	response.Broker = broker
	response.DomainController = controller
	return response
}


func (response *BrokerRegistrationResponse) String() string {
	return fmt.Sprintf("Broker: %s, Controllers: %s",response.Broker, response.DomainController)
}