package domainController

import "github.com/tkrex/IDS/common/models"

type DomainControllerConfiguration struct {
	parentDomain *models.RealWorldDomain
	ownDomain *models.RealWorldDomain

	registrationURL string
}

func NewDomainControllerConfiguration(parentDomain *models.RealWorldDomain, ownDomain *models.RealWorldDomain, registrationURL string) *DomainControllerConfiguration {
	config := new(DomainControllerConfiguration)
	config.parentDomain = parentDomain
	config.ownDomain = ownDomain
	return config
}
