package models

import "fmt"

type BrokerRegistrationResponse struct {
	Broker            *Broker `json:"broker"`
	DomainControllers []*DomainController `json:"domainControllers"`
}

func NewBrokerRegistrationResponse(broker *Broker, controllers []*DomainController) *BrokerRegistrationResponse {
	response := new(BrokerRegistrationResponse)
	response.Broker = broker
	response.DomainControllers = controllers
	return response
}


func (response *BrokerRegistrationResponse) String() string {
	return fmt.Sprintf("Broker: %s, Controllers: %s",response.Broker, response.DomainControllers)
}