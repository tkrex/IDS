package models

type RealWorldDomain struct {
	ID 	int `json:"id"`
	Name 	string `json:"name"`
}

func NewRealWorldDomain(id int, name string) *RealWorldDomain {
	domain := new(RealWorldDomain)
	domain.ID = id
	domain.Name = name
	return domain
}