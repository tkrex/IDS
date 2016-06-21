package registration

import (
	"github.com/tkrex/IDS/common/models"
	"crypto/md5"
	"fmt"
	"errors"
	"github.com/tkrex/IDS/gateway/persistence"
	"github.com/tkrex/IDS/common/routing"
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

	var domainController *models.DomainController
	var err error

	dbWorker := persistence.NewGatewayDBWorker()
	if dbWorker == nil {
		return nil,errors.New("Can't connect with database")
	}
	defer dbWorker.Close()
	domainController, _ = routing.NewRoutingManager().DomainControllerForDomain(broker.RealWorldDomain,false)


	broker.ID = brokerID
	if _, found := dbWorker.FindBrokerById(brokerID); found {
		fmt.Println("Broker Already Registered")
	} else {
		if err = dbWorker.StoreBroker(broker); err != nil {
			return nil, err
		}
	}

	registrationResponse := models.NewBrokerRegistrationResponse(broker, domainController)
	return registrationResponse, nil
}