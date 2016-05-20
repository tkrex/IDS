package common

import (
    "fmt"
	"github.com/tkrex/IDS/common/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)


const (
	Host     = "localhost:27017"
	Username = "example"
	Password = "example"
	Database = "ControlMessages"
	DomainControllerCollection = "domainController"
)



type ControlMessageDBDelegate struct {
	session *mgo.Session
}

func NewControlMessageDBDelegate() (*ControlMessageDBDelegate, error) {
	databaseWorker := new(ControlMessageDBDelegate)
	var error error
	databaseWorker.session, error = mgo.Dial(Host)
	if error != nil {
		return nil , error
	}
	return databaseWorker, error
}

func (dbWoker *ControlMessageDBDelegate)Close() {
	dbWoker.session.Close()
}

func (dbWorker *ControlMessageDBDelegate) domainControllerCollection() *mgo.Collection {
	return dbWorker.session.DB(Database).C(DomainControllerCollection)
}


func (dbWoker *ControlMessageDBDelegate) StoreDomainControllers(domainControllers []*models.DomainController) error {
	coll := dbWoker.domainControllerCollection()
	bulk := coll.Bulk()
	bulk.Unordered()
	for _, domainController := range domainControllers {
		bulk.Upsert(bson.M{"domain.name":domainController.Domain.Name},bson.M{"$set": domainController})
	}
	_, error := bulk.Run()
	return error
}

func (dbWoker *ControlMessageDBDelegate) FindDomainControllerForDomain(domain string) (*models.DomainController,error) {
	coll := dbWoker.domainControllerCollection()
	var domainController *models.DomainController
	error := coll.Find(bson.M{"domain.name":domain}).One(&domainController)
	fmt.Println(error)
	return domainController,error
}

func (worker *ControlMessageDBDelegate) removeDomainControllers(domainControllers []*models.DomainController) error{
	coll := worker.domainControllerCollection()
	bulk := coll.Bulk()
	bulk.Unordered()
	for _, domainController := range domainControllers {
		bulk.Remove(bson.M{"domain.name":domainController.Domain.Name, "ipAddress": domainController.IpAddress})
	}
	_, err := bulk.Run()
	return err
}


