package models

import (
	"time"
)

type Topic struct {
	Name                string        `json:"name" bson:"name"`
	LastPayload         []byte        `json:"payload" bson:"payload"`
	payloadSimilarity   float32        `json:"payloadSimilarity" bson:"payloadSimilarity"`
	LastUpdateTimeStamp time.Time        `json:"lastUpdate" bson:"lastUpdate"`
	UpdateBehavior      *UpdateBehavior
	Domain              *RealWorldDomain   `json:"domain" bson:"domain`
}

type RawTopicMessage struct {
	Name        string
	Payload     []byte
	ArrivalTime time.Time
}

func NewRawTopicMessage(name string, payload []byte) *RawTopicMessage {
	message := new(RawTopicMessage)
	message.Name = name
	message.Payload = payload
	message.ArrivalTime = time.Now()
	return message
}

type UpdateBehavior struct {
	NumberOfUpdates                int                  `json:"numberOfUpdates"`
	UpdateIntervalsInSeconds       []float64           `json:"allIntervals"`
	AverageUpdateIntervalInSeconds float64 `json:"averageInterval"`
	MinimumUpdateIntervalInSeconds int `json:"minimumInterval"`
	MaximumUpdateIntervalInSeconds int `json:"maximumInterval"`
	UpdateReliability              map[string]float64          `json:"reliability"`
}


func NewUpdateBehavior() *UpdateBehavior {
	behavior := new(UpdateBehavior)
	behavior.UpdateReliability = make(map[string]float64)
	return behavior
}

func NewTopic(name string, payload []byte,arrivalTime time.Time) *Topic {
	topic := new(Topic)
	topic.Name = name
	topic.LastPayload = payload
	topic.LastUpdateTimeStamp = arrivalTime
	topic.UpdateBehavior = NewUpdateBehavior()
	return topic
}



