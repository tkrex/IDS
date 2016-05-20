package models

type ControlMessageType string

const (
	DomainControllerUpdate ControlMessageType = "Update"
	DomainControllerDelete ControlMessageType = "Delete"
)

type ControlMessage struct {
	MessageType ControlMessageType `json:"type"`
	DomainControllers []*DomainController `json:"controllers"`
}

func NewControlMessage(messageType ControlMessageType, domainControllers []*DomainController) *ControlMessage {
	controlMessage := new(ControlMessage)
	controlMessage.MessageType = messageType
	controlMessage.DomainControllers = domainControllers
	return controlMessage
}