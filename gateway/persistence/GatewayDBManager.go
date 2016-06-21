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
	DomainCollection = "domains"
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
	fmt.Println("Connected to Database: ",Host)
	worker.session = session
	return worker
}

func (worker *GatewayDBWorker) brokerCollection() *mgo.Collection {
	return worker.session.DB(Database).C(BrokerCollection)
}

func (worker *GatewayDBWorker) domainCollection() *mgo.Collection {
	return worker.session.DB(Database).C(DomainCollection)
}


func (worker *GatewayDBWorker) StoreBroker(broker *models.Broker) (error) {
	coll := worker.brokerCollection()

	if err := coll.Insert(broker); err != nil {
		return err
	}
	return nil
}

func (worker *GatewayDBWorker) StoreDomains(domains []*models.RealWorldDomain) error {
	coll := worker.brokerCollection()
	index := mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	_ = coll.EnsureIndex(index)

	bulkTransaction := coll.Bulk()

	for _, domain := range domains {
		bulkTransaction.Remove(bson.M{"name": domain.Name })
		bulkTransaction.Insert(domain)
	}
	_, err := bulkTransaction.Run()
	return  err
}

func (worker *GatewayDBWorker) FindAllDomains() ([]*models.RealWorldDomain, error) {
	coll := worker.domainCollection()
	domains := []*models.RealWorldDomain{}
	var error error
	if error = coll.Find(nil).All(&domains); error != nil {
		fmt.Println(error)
	}
	return domains, error
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



func (worker *GatewayDBWorker) Close() {
	worker.session.Close()
}