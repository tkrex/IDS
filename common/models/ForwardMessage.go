package models

type ForwardMessage struct {
	Domain   *RealWorldDomain
	Priority int
}

func NewForwardMessage(domain *RealWorldDomain, priority int) *ForwardMessage {
	message := new(ForwardMessage)
	message.Domain = domain
	message.Priority = priority
	return message
}