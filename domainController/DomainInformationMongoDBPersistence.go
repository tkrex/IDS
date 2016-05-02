package domainController

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"github.com/tkrex/IDS/common/models"
	"gopkg.in/mgo.v2/bson"
)

const (
	Host = "localhost:27017"
	Username = "example"
	Password = "example"
	Database = "IDSDomainController"
	Collection = "domainInformation"
)

func OpenSession() *mgo.Session {
	//session, err := mgo.DialWithInfo(&mgo.DialInfo{
	//	Addrs:    []string{Host},
	//	Username: Username,
	//	Password: Password,
	//	Database: Database,
	//	DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
	//		return tls.Dial("tcp", addr.String(), &tls.Config{})
	//	},
	//})

	session, err := mgo.Dial(Host)

	if err != nil {
		panic(err)
	}

	return session
}
func StoreDomainInformation(domainInformationMessage *models.DomainInformationMessage) {

	session := OpenSession()

	defer session.Close()
	fmt.Printf("Connected to %v\n", session.LiveServers())

	coll := session.DB(Database).C(Collection)
	index := mgo.Index{
		Key:        []string{"broker.id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	_ = coll.EnsureIndex(index)

	coll.Remove(bson.M{"broker.id": domainInformationMessage.Broker.ID })
	err := coll.Insert(domainInformationMessage)
	if err != nil {
		fmt.Println(err)
	}
}

func FindAllDomainInformation() ([]*models.DomainInformationMessage, error) {
	session := OpenSession()
	defer session.Close()
	coll := session.DB(Database).C(Collection)
	var domainInformation []*models.DomainInformationMessage
	var error error

	if error := coll.Find(nil).All(&domainInformation); error != nil {
		fmt.Println(error)
	}
	return domainInformation, error
}

func FindDomainInformationByDomainName(domainName string) ([]*models.DomainInformationMessage, error) {
	session := OpenSession()
	defer session.Close()
	var domainInformation []*models.DomainInformationMessage
	var error error
	coll := session.DB(Database).C(Collection)
	if error = coll.Find(bson.M{"realworlddomain.name": domainName}).All(&domainInformation); error != nil {
		fmt.Println(error)
	}
	return domainInformation, error
}

func FindAllBrokers() ([]*models.Broker,error) {
	session := OpenSession()
	defer session.Close()
	coll := session.DB(Database).C(Collection)
	var domainInformation []*models.DomainInformationMessage
	var error error

	if error := coll.Find(nil).Select(bson.M{"broker":1}).All(&domainInformation); error != nil {
		fmt.Printf("Query Error: %s",error)
	}
	brokers := make([]*models.Broker,0,len(domainInformation))
	for _,information := range domainInformation {
		brokers = append(brokers,information.Broker)
	}
	return brokers, error
}

