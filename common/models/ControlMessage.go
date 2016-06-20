package models

type ControlMessageType string

const (
	DomainControllerChange ControlMessageType = "Update"
	DomainControllerDelete ControlMessageType = "Delete"
	DomainControllerFetch ControlMessageType = "Fetch"
)

type ControlMessage struct {
	MessageType      ControlMessageType `json:"type"`
	DomainController *DomainController `json:"controllers"`
}

func NewControlMessage(messageType ControlMessageType, domainController *DomainController) *ControlMessage {
	controlMessage := new(ControlMessage)
	controlMessage.MessageType = messageType
	controlMessage.DomainController = domainController
	return controlMessage
}