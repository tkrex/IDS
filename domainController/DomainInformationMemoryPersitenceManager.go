package domainController

import "github.com/tkrex/IDS/common/models"

type DomainInformationMemoryPersistenceManager struct {
	domainInformation map[int]*models.DomainInformationMessage
}


func NewDomainInformationMemoryPersistenceManager() *DomainInformationMemoryPersistenceManager {
	manager := new(DomainInformationMemoryPersistenceManager)
	manager.domainInformation = make(map[int]*models.DomainInformationMessage)
	return manager
}


func (manager *DomainInformationMemoryPersistenceManager) Store(message *models.DomainInformationMessage) {
	//manager.domainInformation[message.Broker.ID] = message
}


func (manager *DomainInformationMemoryPersistenceManager) DomainInformation() []*models.DomainInformationMessage {
	messageArray := make([]*models.DomainInformationMessage,len(manager.domainInformation))

	index := 0
	for _, message := range manager.domainInformation {
		messageArray[index] = message
		index++
	}
	return messageArray
}

func (manager *DomainInformationMemoryPersistenceManager) DomainInformationWithDomain(domain string) []*models.DomainInformationMessage {
	messageArray := make([]*models.DomainInformationMessage,len(manager.domainInformation))

	index := 0
	for _, message := range manager.domainInformation {
		if message.RealWorldDomain.Name == domain {
			messageArray[index] = message
		}
		index++
	}
	return messageArray
}

func (manager *DomainInformationMemoryPersistenceManager) Brokers() []*models.Broker {
	brokers := make([]*models.Broker,len(manager.domainInformation))

	index := 0
	for _, message := range manager.domainInformation {
		brokers[index] = message.Broker
		index++
	}
	return brokers
}

