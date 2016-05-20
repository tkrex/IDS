package registration

import (
	"github.com/tkrex/IDS/common/models"
	"crypto/md5"
	"fmt"
	"errors"
	"github.com/tkrex/IDS/gateway/persistence"
)

type BrokerRegistrationHandler struct {

}

func NewBrokerRegistrationHandler() *BrokerRegistrationHandler {
	handler := new(BrokerRegistrationHandler)
	return handler
}

func (handler *BrokerRegistrationHandler) RegisterBroker(broker *models.Broker) (*models.BrokerRegistrationResponse, error) {

	brokerIdentificationString := broker.IP + broker.InternetDomain
	byteArray := []byte(brokerIdentificationString)
	md5Bytes := md5.Sum(byteArray)
	brokerID := fmt.Sprintf("%x", md5Bytes)

	var domainControllers []*models.DomainController
	var err error
	dbWorker := persistence.NewGatewayDBWorker()
	if dbWorker == nil {
		return nil,errors.New("Can't connect with database")
	}
	defer dbWorker.Close()
	if domainControllers, err = dbWorker.FindAllDomainController(); err != nil {
		return nil, err
	}
	fmt.Println(domainControllers)

	broker.ID = brokerID
	if _, found := dbWorker.FindBrokerById(brokerID); found {
		fmt.Println("Broker Already Registered")
	} else {
		if err = dbWorker.StoreBroker(broker); err != nil {
			return nil, err
		}
	}

	registrationResponse := models.NewBrokerRegistrationResponse(broker, domainControllers)
	return registrationResponse, nil
}