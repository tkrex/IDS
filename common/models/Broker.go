package models

import (
	"fmt"
	"time"
)

type Broker struct {
	ID string `json:"id"`
	IP string `json:"ip"`
	InternetDomain string `json:"internetDomain"`
	Geolocation *Geolocation `json:"geolocation"`
	RealWorldDomains []*RealWorldDomain `json:"realWorldDomains"`
	Statitics *BrokerStatistic
}

type BrokerStatistic struct {
	NumberOfTopics           int
	ReceivedTopicsPerSeconds float64
	LastStatisticUpdate time.Time
}

func NewBrokerStatistic() *BrokerStatistic {
	statistic := new(BrokerStatistic)
	statistic.NumberOfTopics = 0
	statistic.ReceivedTopicsPerSeconds = 0
	statistic.LastStatisticUpdate = time.Now()
	return statistic
}

func NewBroker() *Broker {
	broker := new(Broker)
	broker.ID = ""
	broker.IP = ""
	broker.InternetDomain = ""
	broker.Geolocation = new(Geolocation)
	broker.RealWorldDomains = []*RealWorldDomain{}
	broker.RealWorldDomains[0] = NewRealWorldDomain("default")
	broker.Statitics = NewBrokerStatistic()
	return broker
}

func (broker *Broker) String() string {
	return fmt.Sprintf("ID: %s, IP: %s, interDomain: %s, geolocation: %s, realWorldDomains: %s", broker.ID,broker.IP,broker.InternetDomain,broker.Geolocation,broker.RealWorldDomains)
}

