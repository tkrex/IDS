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
	Publish(map[string]*models.Topic)
	Close()

}

type InformationProvider interface {

}