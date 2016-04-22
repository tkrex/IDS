package common

import (
	"github.com/tkrex/IDS/common/models"
)

type InformationProducer interface {
	 InformationChannel() chan *models.Topic
	Close()

}

type InformationConsumer interface {
	Close()
}

type InformationProvider interface {

}