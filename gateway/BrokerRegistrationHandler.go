package gateway

import (
	"github.com/tkrex/IDS/common/models"
	"crypto/md5"
	"fmt"
)

type BrokerRegistrationHandler struct {

}

func NewBrokerRegistrationHandler() *BrokerRegistrationHandler {
	handler := new(BrokerRegistrationHandler)
	return handler
}

func (handler *BrokerRegistrationHandler) registerBroker(broker *models.Broker) (*models.BrokerRegistrationResponse, error) {

	brokerIdentificationString := broker.IP + broker.InternetDomain
	byteArray := []byte(brokerIdentificationString)
	md5Bytes := md5.Sum(byteArray)
	brokerID := fmt.Sprintf("%x", md5Bytes)

	var domainControllers []*models.DomainController
	var err error
	if domainControllers, err = FindAllDomainController(); err != nil {
		return nil, err
	}
	fmt.Println(domainControllers)

	broker.ID = brokerID
	if _, found := FindBrokerById(brokerID); found {
		fmt.Println("Broker Already Registered")
	} else {
		if err = StoreBroker(broker); err != nil {
			return nil, err
		}
	}

	registrationResponse := models.NewBrokerRegistrationResponse(broker, domainControllers)
	return registrationResponse, nil
}