package models

import (
	"fmt"
	"strings"
)

type RealWorldDomain struct {
	Name 	string `json:"name"`
}

func NewRealWorldDomain(name string) *RealWorldDomain {
	domain := new(RealWorldDomain)
	domain.Name = name
	return domain
}

func (domain *RealWorldDomain) String() string {
	return fmt.Sprintf("Name: %s", domain.Name)
}

func (domain *RealWorldDomain) FirstLevelDomain() string {
	domainLevels := strings.Split(domain.Name,"/")
	firstLevelDomain := domainLevels[0]
	firstLevelDomain = strings.TrimSpace(firstLevelDomain)
	return firstLevelDomain
}