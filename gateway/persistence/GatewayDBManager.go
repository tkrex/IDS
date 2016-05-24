package persistence

import (
	"gopkg.in/mgo.v2"
	"github.com/tkrex/IDS/common/models"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	Host = "localhost:27017"
	Username = "example"
	Password = "example"
	Database = "IDSGateway"
	BrokerCollection = "brokers"
	DomainControllerCollection = "domainControllers"
	DomainInformationCollection = "domainInformation"
)

type GatewayDBWorker struct {
	session *mgo.Session
}

func NewGatewayDBWorker() *GatewayDBWorker {
	worker := new(GatewayDBWorker)

	session, err := mgo.DialWithTimeout(Host, time.Second * 3)
	if err != nil {
		return nil
	}
	worker.session = session
	return worker
}

func (worker *GatewayDBWorker) brokerCollection() *mgo.Collection {
	return worker.session.DB(Database).C(BrokerCollection)
}

func (worker *GatewayDBWorker) domainControllerCollection() *mgo.Collection {
	return worker.session.DB(Database).C(DomainControllerCollection)
}

func (worker *GatewayDBWorker) StoreBroker(broker *models.Broker) (error) {
	coll := worker.brokerCollection()

	if err := coll.Insert(broker); err != nil {
		return err
	}
	return nil
}

func (worker *GatewayDBWorker) FindAllBrokers() ([]*models.Broker, error) {
	coll := worker.brokerCollection()
	var brokers []*models.Broker
	var error error
	if error = coll.Find(nil).All(brokers); error != nil {
		fmt.Println(error)
	}
	return brokers, error
}

func (worker *GatewayDBWorker) FindBrokerById(brokerID string) (*models.Broker, bool) {
	coll := worker.brokerCollection()

	broker := new(models.Broker)
	if error := coll.Find(bson.M{"id":brokerID}).One(broker); error != nil {
		return nil, false
	}
	return broker, true
}

func (worker *GatewayDBWorker) FindBrokerByIP(brokerIP int) (*models.Broker, error) {
	coll := worker.brokerCollection()

	var broker *models.Broker
	var error error
	if error := coll.Find(bson.M{"ip":brokerIP}).One(broker); error != nil {
		fmt.Println(error)
	}
	return broker, error
}

func (worker *GatewayDBWorker) FindAllDomainController() ([]*models.DomainController, error) {
	var domainControllers []*models.DomainController

	coll := worker.domainControllerCollection()

	if error := coll.Find(nil).All(&domainControllers); error != nil && error != mgo.ErrNotFound {
		return domainControllers, error
	}
	return domainControllers, nil
}


func (worker *GatewayDBWorker) StoreDomainController(domainController *models.DomainController) (bool, error) {
	coll := worker.domainControllerCollection()
	index := mgo.Index{
		Key:        []string{"domain.name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	_ = coll.EnsureIndex(index)
	info, err := coll.Upsert(bson.M{"domain.name":domainController.Domain.Name}, bson.M{"$set": domainController})
	newInformation := info.Updated != 0 || info.Matched == 0
	return newInformation, err
}

func (worker *GatewayDBWorker) RemoveDomainControllerForDomain(domain *models.RealWorldDomain) error {
	coll := worker.domainControllerCollection()
	err := coll.Remove(bson.M{"domain.name":domain.Name})
	return err
}

func (worker *GatewayDBWorker) FindDomainControllerForDomain(domain *models.RealWorldDomain)  *models.DomainController {
	coll := worker.domainControllerCollection()
	var domainController *models.DomainController
	coll.Find(bson.M{"domain.name":domain.Name}).One(domainController)
	return domainController
}

func (worker *GatewayDBWorker) Close() {
	worker.session.Close()
}