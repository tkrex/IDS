package common

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
	Collection = "broker"
)

func openSession() (*mgo.Session, error) {
	session, err := mgo.DialWithTimeout(Host, time.Second * 3)

	if err != nil {
		panic(err)
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

	coll := session.DB(Database).C(Collection)

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
	coll := session.DB(Database).C(Collection)

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
	coll := session.DB(Database).C(Collection)

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
	coll := session.DB(Database).C(Collection)

	var broker *models.Broker
	var error error
	if error := coll.Find(bson.M{"ip":brokerIP}).One(broker); error != nil {
		fmt.Println(error)
	}
	return broker,error
}