package common

import (
	"time"
	"sync"
	"github.com/tkrex/IDS/common/models"
)

type TopicForwarder struct {
	persistenceManager InformationPersistenceManager
	publisher 	InformationPublisher
	forwardInterval time.Duration
	forwardTicker *time.Ticker

	forwarderStarted sync.WaitGroup
	forwarderStopped sync.WaitGroup
}

func NewTopicForwarder(persistenceManager InformationPersistenceManager, forwardInterval time.Duration) *TopicForwarder {
	forwarder := new(TopicForwarder)
	forwarder.persistenceManager = persistenceManager
	forwarder.forwardInterval = forwardInterval
	forwarder.forwarderStarted.Add(1)
	forwarder.forwarderStopped.Add(1)
	go forwarder.run()
	forwarder.forwarderStarted.Wait()
	return forwarder
}

func (forwarder *TopicForwarder) run() {
	forwarder.forwardTicker = time.NewTicker(forwarder.forwardInterval)
	config := models.NewMqttClientConfiguration("tcp://localhost:1883","burst","publisher")
	forwarder.publisher = NewMqttPublisher(config)
	forwarder.forwarderStarted.Done()
	go func() {
		for _ = range forwarder.forwardTicker.C {
			topics := forwarder.persistenceManager.Topics()
			go forwarder.publisher.PublishTopics(topics)
		}
	}()
}
