package common

import (
	"github.com/tkrex/IDS/common/models"
)

type InformationProducer interface {
	Run()
	Close()
}

type InformationProcessor interface {
	Close()
}

type InformationPublisher interface {
	PublishTopics([]models.Topic) error
	Close()

}

type InformationPersistenceManager interface {
	Topics() [] models.Topic
	TopicWithName(string) ( models.Topic,bool)
	StoreTopic( models.Topic)
	NumberOfTopics() int
}

type InformationProvider interface {

}