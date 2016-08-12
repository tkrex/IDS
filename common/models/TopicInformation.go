package models

import (
	"time"
)

//Contains all information collected for a Topic
type TopicInformation struct {
	Name                    string        `json:"name" bson:"name"`
	LastPayload             string       `json:"payload" bson:"payload"`
	PayloadSimilarity       float64       `json:"payloadSimilarity" bson:"payloadSimilarity"`
	SimilarityCheckInterval int                `json:"-" bson:""`
	FirstUpdateTimeStamp    time.Time            `json:"firstUpdate"`
	LastUpdateTimeStamp     time.Time        `json:"lastUpdate" bson:"lastUpdate"`
	UpdateBehavior          *UpdateBehavior	   `json:"updateBehavior" bson:"updateBehavior`
	Domain                  *RealWorldDomain   `json:"domain" bson:"domain`
	Visibility              bool                `json:"-"`
}

func NewTopicInformation(name string, payload string, arrivalTime time.Time) *TopicInformation {
	topic := new(TopicInformation)
	topic.Name = name
	topic.LastPayload = payload
	topic.FirstUpdateTimeStamp = arrivalTime
	topic.LastUpdateTimeStamp = arrivalTime
	topic.SimilarityCheckInterval = 1
	topic.Visibility = true
	topic.UpdateBehavior = NewUpdateBehavior()
	return topic
}

//Contains the information of a MQTT messages, as well as the arrival time
type RawTopicMessage struct {
	Topic       string
	Payload     []byte
	ArrivalTime time.Time
}

func NewRawTopicMessage(topic string, payload []byte) *RawTopicMessage {
	message := new(RawTopicMessage)
	message.Topic = topic
	message.Payload = payload
	message.ArrivalTime = time.Now()
	return message
}

//Contains information about the time-related behavior of a Topic
type UpdateBehavior struct {
	NumberOfUpdates                int                  `json:"numberOfUpdates"`
	UpdateIntervalsInSeconds       []float64           `json:"allIntervals"`
	AverageUpdateIntervalInSeconds float64 `json:"averageInterval"`
	MinimumUpdateIntervalInSeconds int `json:"minimumInterval"`
	MaximumUpdateIntervalInSeconds int `json:"maximumInterval"`
	UpdateIntervalDeviation        float64          `json:"deviation"`
	Reliability                    string `json:"reliability"`
}

func NewUpdateBehavior() *UpdateBehavior {
	behavior := new(UpdateBehavior)
	return behavior
}





