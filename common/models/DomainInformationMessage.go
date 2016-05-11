package models

type DomainInformationMessage struct {
	RealWorldDomain *RealWorldDomain `json:"domain" bson:"domain"`
	ForwardPriority int `json:"forwardPriority"`
	Broker   *Broker `json:"broker" bson:"broker"`
	Topics   []*Topic `json:"topics" bson:"topics"`
}

func NewDomainInformationMessage(domain *RealWorldDomain, broker *Broker, topics []*Topic) *DomainInformationMessage {
	message := new(DomainInformationMessage)
	message.RealWorldDomain = domain
	message.Topics = topics
	message.Broker = broker
	return message
}
