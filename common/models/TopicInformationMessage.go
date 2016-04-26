package models

type TopicInformationMessage struct {
	Broker *Broker `json:"broker"`
	Topcis []Topic `json:"topics"`
}

func NewTopicInformationMessage(broker *Broker, topics []Topic) *TopicInformationMessage{
	message := new(TopicInformationMessage)
	message.Topcis = topics
	message.Broker = broker
	return message
}
