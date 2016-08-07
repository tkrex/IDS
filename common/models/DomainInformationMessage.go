package models

import "fmt"

type DomainInformationMessage struct {
	RealWorldDomain *RealWorldDomain `json:"domain" bson:"domain"`
	ForwardPriority int `json:"forwardPriority"`
	Broker   *Broker `json:"broker" bson:"broker"`
	Topics   []*TopicInformation `json:"topics" bson:"topics"`
}

func NewDomainInformationMessage(domain *RealWorldDomain, broker *Broker, topics []*TopicInformation) *DomainInformationMessage {
	message := new(DomainInformationMessage)
	message.RealWorldDomain = domain
	message.Topics = topics
	message.Broker = broker
	return message
}

func (message *DomainInformationMessage) String() string {
	return fmt.Sprintf("Broker: %s, Domain: %s, Topics: %s, ForwardPriority: %d",message.Broker,message.RealWorldDomain,message.Topics,message.ForwardPriority)
}
