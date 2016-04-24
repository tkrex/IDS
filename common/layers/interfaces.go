package common

import (
	"github.com/tkrex/IDS/common/models"
)

type InformationProducer interface {
	InformationChannel() chan *models.Topic
	Close()

}

type InformationProcessor interface {
	Close()
}

type InformationPublisher interface {
	PublishTopics([]*models.Topic)
	Close()

}

type InformationPersistenceManager interface {
	Topics() []*models.Topic
	TopicWithName(string) (*models.Topic,bool)
	StoreTopic(*models.Topic)
}

type InformationProvider interface {

}