package models

import (
	"time"
	"github.com/tkrex/IDS/common"
)


type Topic struct {
	ID		    int        `json:"id"`
	Name                string        `json:"name"`
	LastPayload         []byte        `json:"payload"`
	payloadSimilarity   float32        `json:"payloadSimilarity"`
	LastUpdateTimeStamp time.Time        `json:"lastUpdate"`
	UpdateBehavior      *UpdateBehavior
	Domain		RealWorldDomain   `json:"domain"`
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

func NewTopic(id int, name string, payload []byte) *Topic {
	topic := new(Topic)
	topic.ID = id
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

