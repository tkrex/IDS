package configuration

import (
	"gopkg.in/mgo.v2"
	"os"
)

const (
	Host = "localhost:27017"
	Username = "example"
	Password = "example"
	Database = "IDSDomainController"

	DomainControllerConfigCollection = "domainControllerConfig"
)

type DomainControllerConfigStorageManager struct {
	session *mgo.Session
}

func NewDomainControllerConfigStorageManager() (*DomainControllerConfigStorageManager, error) {
	databaseWorker := new(DomainControllerConfigStorageManager)
	var error error
	databaseWorker.session, error = openSession()
	if error != nil {
		return nil, error
	}
	return databaseWorker, error
}

func openSession() (*mgo.Session, error) {
	host := os.Getenv("MONGODB_URI")
	session, err := mgo.Dial(host)
	return session, err
}

func (dbWoker *DomainControllerConfigStorageManager)Close() {
	dbWoker.session.Close()
}

func (dbWorker *DomainControllerConfigStorageManager) domainControllerConfigCollection() *mgo.Collection {
	return dbWorker.session.DB(Database).C(DomainControllerConfigCollection)
}


func (dbWoker *DomainControllerConfigStorageManager) StoreDomainControllerConfig(domainControllerConfig *DomainControllerConfiguration) error {
	configCollection := dbWoker.domainControllerConfigCollection()
	err := configCollection.Insert(domainControllerConfig)
	return err
}

func (dbWoker *DomainControllerConfigStorageManager) FindDomainControllerConfig() (*DomainControllerConfiguration,error) {
	configCollection := dbWoker.domainControllerConfigCollection()
	var config *DomainControllerConfiguration
	err := configCollection.Find(nil).One(&config)
	return config, err
}
