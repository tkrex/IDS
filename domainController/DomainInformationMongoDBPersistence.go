package domainController

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"github.com/tkrex/IDS/common/models"
	"gopkg.in/mgo.v2/bson"
	"github.com/tkrex/IDS/common"
)


const (
	Host     = "localhost:27017"
	Username = "example"
	Password = "example"
	Database = "IDSDomainController"

	BrokerCollection = "brokers"
	DomainInformationCollection = "domainInformation"
)



type DomainControllerDatabaseWorker struct {
	session *mgo.Session
}

func NewDomainControllerDatabaseWorker() (*DomainControllerDatabaseWorker, error) {
	databaseWorker := new(DomainControllerDatabaseWorker)
	var error error
	databaseWorker.session, error = openSession()
	if error != nil {
		return nil , error
	}
	return databaseWorker, error
}


func openSession() (*mgo.Session,error) {
	session , err := mgo.Dial(Host)
	return session, err
}

func (dbWoker *DomainControllerDatabaseWorker)Close() {
	dbWoker.session.Close()
}


func (dbWorker *DomainControllerDatabaseWorker) domainInformationCollection() *mgo.Collection {
	return dbWorker.session.DB(Database).C(DomainInformationCollection)
}


func (dbWorker *DomainControllerDatabaseWorker) brokerCollection() *mgo.Collection {
	return dbWorker.session.DB(Database).C(BrokerCollection)
}


type DomainQuery struct {
	Domain *models.RealWorldDomain `bson:"domain"`
}

func (dbWorker* DomainControllerDatabaseWorker) FindAllDomains() ([]*models.RealWorldDomain,error) {
	coll := dbWorker.domainInformationCollection()


	domainFields := []*DomainQuery{}
	domains := []*models.RealWorldDomain{}
	err := coll.Find(nil).Select(bson.M{"_id": 0, "domain": 1}).All(&domainFields)
	for _, domainQuery := range domainFields {
		domain := domainQuery.Domain
		if !common.Include(domains,domain) {
			domains = append(domains,domain)
		}
	}
	return domains, err
}

func (dbWoker *DomainControllerDatabaseWorker) StoreBroker(broker *models.Broker) (error) {
	coll := dbWoker.brokerCollection()

	if err := coll.Insert(broker); err != nil {
		return err
	}
	return nil
}

func (dbWoker *DomainControllerDatabaseWorker) FindBroker() (*models.Broker,error) {
	coll := dbWoker.brokerCollection()

	var error error
	broker := new(models.Broker)
	if error = coll.Find(nil).One(broker); error != nil {
		fmt.Println(error)
	}
	return broker, error
}


func (dbWoker *DomainControllerDatabaseWorker) StoreDomainInformation(domainInformationMessages []*models.DomainInformationMessage) error {
	coll := dbWoker.domainInformationCollection()
	index := mgo.Index{
		Key:        []string{"broker.id","domain.name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	_ = coll.EnsureIndex(index)

	dbWoker.RemoveDomainInformation(domainInformationMessages)
	transaction := coll.Bulk()

	for _,information := range domainInformationMessages {
		transaction.Insert(information)
	}
	_,err := transaction.Run()
	return err
}

func (dbWoker *DomainControllerDatabaseWorker) RemoveDomainInformation(domainInformationMessages []*models.DomainInformationMessage) error {
	coll := dbWoker.domainInformationCollection()
	transaction := coll.Bulk()
	transaction.Unordered()
	for _,information := range domainInformationMessages {
		transaction.Remove(bson.M{"broker.id": information.Broker.ID, "domain.name": information.RealWorldDomain.Name })
	}
	_, err := transaction.Run()
	return err
}



func (dbWoker *DomainControllerDatabaseWorker) FindAllDomainInformation() ([]*models.DomainInformationMessage, error) {

	coll := dbWoker.domainInformationCollection()
	var domainInformation []*models.DomainInformationMessage
	var error error

	if error := coll.Find(nil).All(&domainInformation); error != nil {
		fmt.Println(error)
	}
	return domainInformation, error
}

func (dbWoker *DomainControllerDatabaseWorker) FindDomainInformationByDomainName(domainName string) ([]*models.DomainInformationMessage, error) {
	var domainInformation []*models.DomainInformationMessage
	var error error
	coll := dbWoker.domainInformationCollection()
	if error = coll.Find(bson.M{"domain.name": domainName}).All(&domainInformation); error != nil {
		fmt.Println(error)
	}
	return domainInformation, error
}

func (dbWoker *DomainControllerDatabaseWorker) FindAllBrokers() ([]*models.Broker,error) {
	coll := dbWoker.brokerCollection()
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

