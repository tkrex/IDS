package models

type DomainInformationMessage struct {
	RealWorldDomain *RealWorldDomain `json:"domain"`
	Broker   *Broker `json:"broker"`
	Topics   []*Topic `json:"topics"`
}

func NewDomainInformationMessage(domain *RealWorldDomain, broker *Broker, topics []*Topic) *DomainInformationMessage {
	message := new(DomainInformationMessage)
	message.RealWorldDomain = domain
	message.Topics = topics
	message.Broker = broker
	return message
}
