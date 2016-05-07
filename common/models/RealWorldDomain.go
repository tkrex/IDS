package models

import "fmt"

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