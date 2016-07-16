package controlling

import (
	"github.com/tkrex/IDS/common/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"fmt"
	"os"
)

const (
	Host = "localhost:27017"
	Username = "example"
	Password = "example"
	Database = "ControlMessages"
	DomainControllerCollection = "domainController"
)

type DomainControllerStorageManager struct {
	session *mgo.Session
}

func NewDomainControllerStorageManager() *DomainControllerStorageManager {
	databaseWorker := new(DomainControllerStorageManager)
	var error error
	databaseWorker.session, error = mgo.Dial(Host)
	if error != nil {
		fmt.Println("Cant Connect to Database")
		os.Exit(1)
	}
	return databaseWorker
}

func (dbWoker *DomainControllerStorageManager)Close() {
	dbWoker.session.Close()
}

func (dbWorker *DomainControllerStorageManager) domainControllerCollection() *mgo.Collection {
	return dbWorker.session.DB(Database).C(DomainControllerCollection)
}

func (worker *DomainControllerStorageManager) FindAllDomainController() ([]*models.DomainController, error) {
	var domainControllers []*models.DomainController

	coll := worker.domainControllerCollection()

	if error := coll.Find(nil).All(&domainControllers); error != nil && error != mgo.ErrNotFound {
		return domainControllers, error
	}
	return domainControllers, nil
}

func (worker *DomainControllerStorageManager) StoreDomainController(domainController *models.DomainController) (bool, error) {
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

func (worker *DomainControllerStorageManager) RemoveDomainControllerForDomain(domain *models.RealWorldDomain) error {
	coll := worker.domainControllerCollection()
	err := coll.Remove(bson.M{"domain.name":domain.Name})
	return err
}

func (worker *DomainControllerStorageManager) FindDomainControllerForDomain(domain *models.RealWorldDomain) *models.DomainController {
	coll := worker.domainControllerCollection()
	var domainController *models.DomainController
	if err := coll.Find(bson.M{"domain.name":domain.Name}).One(&domainController); err != nil {
		fmt.Println(err)
	}
	return domainController
}


