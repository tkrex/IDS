package models

import (
	"time"
	"github.com/tkrex/IDS/common"
)


type Topic struct {
	Name                string        `json:"name" bson:"name"`
	LastPayload         []byte        `json:"payload" bson:"payload"`
	payloadSimilarity   float32        `json:"payloadSimilarity" bson:"payloadSimilarity"`
	LastUpdateTimeStamp time.Time        `json:"lastUpdate" bson:"lastUpdate"`
	UpdateBehavior      *UpdateBehavior
	Domain		RealWorldDomain   `json:"domain" bson:"domain`
}


type RawTopicMessage struct {
	Name string
	Payload []byte
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
	NumberOfUpdates	    int		  `json:"numberOfUpdates"`
	AverageUpdateIntervalInSeconds int `json:"averageInterval"`
	MinimumUpdateIntervalInSeconds int `json:"minimumInterval"`
	MaximumUpdateIntervalInSeconds int `json:"maximumInterval"`
	UpdateReliability float32          `json:"reliability"`
}


func NewUpdateBehavior() *UpdateBehavior {
	behavior := new(UpdateBehavior)
	return behavior
}

func NewTopic(name string, payload []byte) *Topic {
	topic := new(Topic)
	topic.Name = name
	topic.LastPayload = payload
	topic.LastUpdateTimeStamp = time.Now()
	topic.UpdateBehavior = NewUpdateBehavior()
	return topic
}

func (topic *Topic) CalculateUpdateBehavior(newUpdateInterval int) {
	if topic.UpdateBehavior.NumberOfUpdates == 0 {
		topic.UpdateBehavior.NumberOfUpdates++
		return
	}
	if topic.UpdateBehavior.NumberOfUpdates == 1 {
		topic.UpdateBehavior.AverageUpdateIntervalInSeconds = newUpdateInterval
		topic.UpdateBehavior.MaximumUpdateIntervalInSeconds = newUpdateInterval
		topic.UpdateBehavior.MinimumUpdateIntervalInSeconds = newUpdateInterval
	} else {
		topic.UpdateBehavior.MaximumUpdateIntervalInSeconds = common.Max(topic.UpdateBehavior.MaximumUpdateIntervalInSeconds,newUpdateInterval)
		topic.UpdateBehavior.MinimumUpdateIntervalInSeconds = common.Min(topic.UpdateBehavior.MinimumUpdateIntervalInSeconds,newUpdateInterval)
		topic.UpdateBehavior.AverageUpdateIntervalInSeconds = (topic.UpdateBehavior.AverageUpdateIntervalInSeconds * topic.UpdateBehavior.NumberOfUpdates + newUpdateInterval) / (topic.UpdateBehavior.NumberOfUpdates + 1)
	}
	topic.UpdateBehavior.NumberOfUpdates++
}

