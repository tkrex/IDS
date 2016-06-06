package controlling

import (
	"github.com/tkrex/IDS/common/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	Host = "localhost:27017"
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
		return nil, error
	}
	return databaseWorker, error
}

func (dbWoker *ControlMessageDBDelegate)Close() {
	dbWoker.session.Close()
}

func (dbWorker *ControlMessageDBDelegate) domainControllerCollection() *mgo.Collection {
	return dbWorker.session.DB(Database).C(DomainControllerCollection)
}

func (worker *ControlMessageDBDelegate) FindAllDomainController() ([]*models.DomainController, error) {
	var domainControllers []*models.DomainController

	coll := worker.domainControllerCollection()

	if error := coll.Find(nil).All(&domainControllers); error != nil && error != mgo.ErrNotFound {
		return domainControllers, error
	}
	return domainControllers, nil
}


func (worker *ControlMessageDBDelegate) StoreDomainController(domainController *models.DomainController) (bool, error) {
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

func (worker *ControlMessageDBDelegate) RemoveDomainControllerForDomain(domain *models.RealWorldDomain) error {
	coll := worker.domainControllerCollection()
	err := coll.Remove(bson.M{"domain.name":domain.Name})
	return err
}

func (worker *ControlMessageDBDelegate) FindDomainControllerForDomain(domainName string)  *models.DomainController {
	coll := worker.domainControllerCollection()
	var domainController *models.DomainController
	coll.Find(bson.M{"domain.name":domainName}).One(domainController)
	return domainController
}


