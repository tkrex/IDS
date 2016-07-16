package models

import (
	"fmt"
	"strings"
)

type RealWorldDomain struct {
	Name string `json:"name"`
}

func NewRealWorldDomain(name string) *RealWorldDomain {
	domain := new(RealWorldDomain)
	domain.Name = name
	return domain
}

func (domain *RealWorldDomain) String() string {
	return fmt.Sprintf("Name: %s", domain.Name)
}

func (domain *RealWorldDomain) TopLevelDomain() *RealWorldDomain {
	domainLevels := domain.DomainLevels()
	return NewRealWorldDomain(domainLevels[0])
}

func (domain *RealWorldDomain) DomainLevels() []string {
	domainLevels := strings.Split(domain.Name, "/")
	for index, domain := range domainLevels {
		domainLevels[index] = strings.TrimSpace(domain)
	}
	return domainLevels
}

func (domain *RealWorldDomain) IsSubDomainOf(secondDomain *RealWorldDomain) bool {
	firstDomainLevels := domain.DomainLevels()
	secondDomainLevels := secondDomain.DomainLevels()
	if len(firstDomainLevels) <= len(secondDomainLevels) {
		return false
	}

	for index, domain := range secondDomainLevels {
		if domain != firstDomainLevels[index] {
			return false
		}
	}
	return true
}

func (domain *RealWorldDomain) ParentDomain() *RealWorldDomain {
	domainLevels := domain.DomainLevels()
	if len(domainLevels) <= 1 {
		return nil
	}
	domainString := ""
	for i := 0;i < len(domainLevels) - 1; i++{
		domainString += domainLevels[i]
		if i < len(domainLevels) -2 {
			domainString += "/"
		}
	}
	return NewRealWorldDomain(domainString)
}