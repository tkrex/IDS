package common

import (
	"github.com/tkrex/IDS/common/models"
	"errors"
	"crypto/md5"
	"fmt"
)

type BrokerRegistrationHandler struct {

}

func NewBrokerRegistrationHandler() *BrokerRegistrationHandler {
	handler := new(BrokerRegistrationHandler)
	return handler
}

func (handler *BrokerRegistrationHandler) registerBroker(broker *models.Broker) (*models.Broker, error) {

	brokerIdentificationString := broker.IP + broker.InternetDomain
	byteArray := []byte(brokerIdentificationString)
	md5Bytes := md5.Sum(byteArray)
	brokerID := fmt.Sprintf("%x", md5Bytes)

	if _,found := FindBrokerById(brokerID); found {
		return nil, errors.New("Broker Already Registered")
	}

	broker.ID = brokerID
	if err := StoreBroker(broker); err != nil {
		return nil, err
	}
	return broker, nil
}