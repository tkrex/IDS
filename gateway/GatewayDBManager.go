package gateway

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
	DomainControllerConnection = "domainControllers"
)

func openSession() (*mgo.Session, error) {
	session, err := mgo.DialWithTimeout(Host, time.Second * 3)

	if err != nil {
		return nil,err
	}
	return session, nil
}

func StoreBroker(broker *models.Broker) (error) {
	session, error :=  openSession()
	defer session.Close()
	if error != nil {
		return error
	}
	fmt.Printf("Connected to %v\n", session.LiveServers())

	coll := session.DB(Database).C(BrokerCollection)

	if err := coll.Insert(broker); err != nil {
		return err
	}
	return nil
}

func FindAllBrokers() ([]*models.Broker,error) {
	var brokers []*models.Broker

	session,err :=  openSession()
	if err != nil {
		return brokers,err
	}
	defer session.Close()
	coll := session.DB(Database).C(BrokerCollection)

	var error error
	if error = coll.Find(nil).All(brokers); error != nil {
		fmt.Println(error)
	}
	return brokers, error
}

func FindBrokerById(brokerID string) (*models.Broker,bool) {
	session, err :=  openSession()
	if err != nil {
		return nil,false
	}
	defer session.Close()
	coll := session.DB(Database).C(BrokerCollection)

	broker := new(models.Broker)
	if error := coll.Find(bson.M{"id":brokerID}).One(broker); error != nil {
		return nil, false
	}
	return broker,true
}

func FindBrokerByIP(brokerIP int) (*models.Broker,error) {
	session, err :=  openSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	coll := session.DB(Database).C(BrokerCollection)

	var broker *models.Broker
	var error error
	if error := coll.Find(bson.M{"ip":brokerIP}).One(broker); error != nil {
		fmt.Println(error)
	}
	return broker,error
}

func FindAllDomainController() ([]*models.DomainController,error) {
	var domainControllers []*models.DomainController

	session,err :=  openSession()
	if err != nil {
		return domainControllers,err
	}
	defer session.Close()
	coll := session.DB(Database).C(DomainControllerConnection)

	if error := coll.Find(nil).All(&domainControllers); error != nil && error != mgo.ErrNotFound {
		return domainControllers,error
	}
	return domainControllers, nil
}



//true: entry was updated
//false: no new data
func UpdateControllerInformation(domainController *models.DomainController) (bool, error) {
	info, error := storeDomainController(domainController)
	return info.Updated != 0, error
}

func storeDomainController(domainController *models.DomainController) (*mgo.ChangeInfo, error) {
	session,err :=  openSession()
	if err != nil {
		return nil,err
	}
	defer session.Close()
	coll := session.DB(Database).C(DomainControllerConnection)
	index := mgo.Index{
		Key:        []string{"domain.name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	_ = coll.EnsureIndex(index)
	info,err := coll.Upsert(bson.M{"domain.name":domainController.Domain.Name},bson.M{"$set": domainController})
	return info, err
}