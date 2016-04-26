package models

type Broker struct {
	ID int `json:"id"`
	IP string `json:"ip"`
	InternetDomain string `json:"internetDomain"`
	Geolocation *Geolocation `json:"geolocation"`
}

func NewBroker(id int, ip string, interDomain string) *Broker {
	broker := new(Broker)
	broker.ID = id
	broker.IP = ip
	broker.InternetDomain = interDomain
	broker.Geolocation = NewGeolocation("Deutschland","Bayern", "München")
	return broker
}

