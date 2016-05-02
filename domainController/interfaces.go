package domainController

import "github.com/tkrex/IDS/common/models"

type DomainInformationPersistenceManager interface {
	Store(*models.DomainInformationMessage)
	DomainInformation() []*models.DomainInformationMessage
	DomainInformationWithDomain(domain string) []*models.DomainInformationMessage
	Brokers() []*models.Broker
}
