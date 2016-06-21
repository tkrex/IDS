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

func (domain *RealWorldDomain) FirstLevelDomain() *RealWorldDomain {
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

	if len(firstDomainLevels) < len(secondDomainLevels) {
		return false
	}

	for index, domain := range secondDomainLevels {
		if domain != firstDomainLevels[index] {
			return false
		}
	}
	return true
}