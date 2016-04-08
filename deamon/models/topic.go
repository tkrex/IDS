package models

type Topic struct {
	id string
	name string
	broker Broker
	lastPayload []byte
	lastUpdateTimeStamp int64
	updateInterval int64
} 
