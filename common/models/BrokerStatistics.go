package models

import "time"

//Contains Broker Statistics
type BrokerStatistics struct {
NumberOfTopics           int `json:"numberOfTopics"`
ReceivedTopicsPerSeconds float64 `json:"topicsPerSecond"`
LastStatisticUpdate      time.Time `json:"-"`
}

func NewBrokerStatistics() *BrokerStatistics {
statistic := new(BrokerStatistics)
statistic.NumberOfTopics = 0
statistic.ReceivedTopicsPerSeconds = 0
statistic.LastStatisticUpdate = time.Now()
return statistic
}