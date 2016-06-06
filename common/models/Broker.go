package models

import (
	"fmt"
	"time"
)

type Broker struct {
	ID              string `json:"id"`
	IP              string `json:"ip"`
	InternetDomain  string `json:"internetDomain"`
	Geolocation     *Geolocation `json:"geolocation"`
	RealWorldDomain *RealWorldDomain `json:"realWorldDomain"`
	Statitics       *BrokerStatistic `json:"statitics"`
}

type BrokerStatistic struct {
	NumberOfTopics           int `json:"numberOfTopics"`
	ReceivedTopicsPerSeconds float64 `json:"topicsPerSecond"`
	LastStatisticUpdate      time.Time `json:"-"`
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
	broker.RealWorldDomain = NewRealWorldDomain("default")
	broker.Statitics = NewBrokerStatistic()
	return broker
}

func (broker *Broker) String() string {
	return fmt.Sprintf("ID: %s, IP: %s, interDomain: %s, geolocation: %s, realWorldDomains: %s", broker.ID, broker.IP, broker.InternetDomain, broker.Geolocation, broker.RealWorldDomain)
}

