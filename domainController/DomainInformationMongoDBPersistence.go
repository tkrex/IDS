package domainController

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"github.com/tkrex/IDS/common/models"
	"gopkg.in/mgo.v2/bson"
)

const (
	Host     = "localhost:27017"
	Username = "example"
	Password = "example"
	Database = "IDSDomainController"
	Collection = "domainInformarmation"
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

	session , err := mgo.Dial(Host)

	if err != nil {
		panic(err)
	}

	return session
}
func  StoreDomainInformation(domainInformationMessages []*models.DomainInformationMessage) {

	session := OpenSession()

	defer session.Close()
	fmt.Printf("Connected to %v\n", session.LiveServers())

	coll := session.DB(Database).C(Collection)
	index := mgo.Index{
		Key:        []string{"broker.id","domain.id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	_ = coll.EnsureIndex(index)


	bulkTransaction := coll.Bulk()

	for _,message := range domainInformationMessages {
		bulkTransaction.Remove(bson.M{"broker.id": message.Broker.ID })
		bulkTransaction.Insert(message)

	}
	bulkResult, err := bulkTransaction.Run()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(bulkResult)
}


func FindAllDomainInformation() ([]*models.DomainInformationMessage, error) {
	session := OpenSession()
	defer session.Close()
	coll := session.DB(Database).C(Collection)
	var domainInformation []*models.DomainInformationMessage
	if err := coll.Find(nil).All(&domainInformation); err != nil {
		fmt.Println(err)
		return err
	}
	return domainInformation
}

func FindDomainInformationByDomainName(domainName string) (*models.DomainInformationMessage, error) {
	session := OpenSession()
	defer session.Close()
	var domainInformation *models.DomainInformationMessage
	coll := session.DB(Database).C(Collection)
	if err := coll.Find(bson.M{"domain.name": domainName}).One(&domainInformation); err != nil {
		fmt.Println(err)
		return err
	}
	return domainInformation
}

