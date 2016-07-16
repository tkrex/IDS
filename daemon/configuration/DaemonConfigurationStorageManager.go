package configuration

import (
	"gopkg.in/mgo.v2"
)

const (
	Database = "IDSDaemon"
	DaemonConfigCollection = "daemonConfig"
)

type DaemonConfigStorageManager struct {
	session *mgo.Session
}

func NewDaemonConfigStorageManager() (*DaemonConfigStorageManager, error) {
	databaseWorker := new(DaemonConfigStorageManager)
	var error error
	databaseWorker.session, error = openSession()
	if error != nil {
		return nil, error
	}
	return databaseWorker, error
}

func openSession() (*mgo.Session, error) {
	Host := DaemonConfigurationManagerInstance().Config().DatabaseAddress
	session, err := mgo.Dial(Host)
	return session, err
}

func (dbWoker *DaemonConfigStorageManager)Close() {
	dbWoker.session.Close()
}

func (dbWorker *DaemonConfigStorageManager) daemonConfigCollection() *mgo.Collection {
	return dbWorker.session.DB(Database).C(DaemonConfigCollection)
}


func (dbWoker *DaemonConfigStorageManager) StoreDaemonConfig(daemonConfig *DaemonConfiguration) error {
	configCollection := dbWoker.daemonConfigCollection()
	err := configCollection.Insert(daemonConfig)
	return err
}

func (dbWoker *DaemonConfigStorageManager) FindDaemonConfig() (*DaemonConfiguration,error) {
	configCollection := dbWoker.daemonConfigCollection()
	var config *DaemonConfiguration
	err := configCollection.Find(nil).One(&config)
	return config, err
}