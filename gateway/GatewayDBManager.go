package gateway

import (
	"gopkg.in/mgo.v2"
	"github.com/tkrex/IDS/common/models"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"time"
	"github.com/tkrex/IDS/domainController"
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

//true: entry was updated
//false: no new data
func (worker *GatewayDBWorker) UpdateControllerInformation(domainController *models.DomainController) (bool, error) {
	info, error := worker.storeDomainController(domainController)
	newInformation := info.Updated != 0 || info.Matched == 0
	return newInformation, error
}

func (worker *GatewayDBWorker) storeDomainController(domainController *models.DomainController) (*mgo.ChangeInfo, error) {
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
	return info, err
}

func (worker *GatewayDBWorker) removeDomainControllers(domainControllers []*models.DomainController) error{
	coll := worker.domainControllerCollection()
	bulk := coll.Bulk()
	bulk.Unordered()
	for _, domainController := range domainControllers {
		bulk.Remove(bson.M{"domain.name":domainController.Domain.Name, "ipAddress": domainController.IpAddress})
	}
	_, err := bulk.Run()
	return err
}

func (worker *GatewayDBWorker) Close() {
	worker.session.Close()
}