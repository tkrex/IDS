package models

type RealWorldDomain struct {
	Name 	string `json:"name"`
}

func NewRealWorldDomain(name string) *RealWorldDomain {
	domain := new(RealWorldDomain)
	domain.Name = name
	return domain
}