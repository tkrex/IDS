package models

import "fmt"

type Broker struct {
	ID string `json:"id"`
	IP string `json:"ip"`
	InternetDomain string `json:"internetDomain"`
	Geolocation *Geolocation `json:"geolocation"`
	RealWorldDomains []*RealWorldDomain `json:"realWorldDomains"`
}

func NewBroker(ip string, interDomain string) *Broker {
	broker := new(Broker)
	broker.ID = ""
	broker.IP = ip
	broker.InternetDomain = interDomain
	broker.Geolocation = new(Geolocation)
	broker.RealWorldDomains = []*RealWorldDomain{}
	return broker
}

func (broker *Broker) String() string {
	return fmt.Sprintf("ID: %s, IP: %s, interDomain: %s, geolocation: %s, realWorldDomains: %s", broker.ID,broker.IP,broker.InternetDomain,broker.Geolocation,broker.RealWorldDomains)
}

