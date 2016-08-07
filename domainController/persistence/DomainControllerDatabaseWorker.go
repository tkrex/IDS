package persistence

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"github.com/tkrex/IDS/common/models"
	"gopkg.in/mgo.v2/bson"
)

const (
	Database = "IDSDomainController"
	BrokerCollection = "brokers"
	DomainInformationCollection = "domainInformation"
	DomainCollection = "domains"
)

type DomainInformationStorage struct {
	session *mgo.Session
}

func NewDomainInformationStorage() (*DomainInformationStorage, error) {
	databaseWorker := new(DomainInformationStorage)
	var error error
	databaseWorker.session, error = openSession()
	if error != nil {
		return nil, error
	}
	return databaseWorker, error
}

func openSession() (*mgo.Session, error) {
	//host := os.Getenv("MONGODB_URI")
	host := "db-default:27017"
	session, err := mgo.Dial(host)
	return session, err
}

func (dbWoker *DomainInformationStorage)Close() {
	dbWoker.session.Close()
}

func (dbWorker *DomainInformationStorage) domainInformationCollection() *mgo.Collection {
	return dbWorker.session.DB(Database).C(DomainInformationCollection)
}

func (dbWorker *DomainInformationStorage) domainCollection() *mgo.Collection {
	return dbWorker.session.DB(Database).C(DomainCollection)
}

func (dbWorker *DomainInformationStorage) brokerCollection() *mgo.Collection {
	return dbWorker.session.DB(Database).C(BrokerCollection)
}

func (dbWoker *DomainInformationStorage) StoreDomainInformationMessage(domainInformationMessage *models.DomainInformationMessage) error {
	domainInformationCollection := dbWoker.domainInformationCollection()
	index := mgo.Index{
		Key:        []string{"broker.id", "domain.name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	_ = domainInformationCollection.EnsureIndex(index)

	err := domainInformationCollection.Insert(domainInformationMessage)
	err = dbWoker.StoreDomain(domainInformationMessage.RealWorldDomain)
	return err
}

func (dbWoker *DomainInformationStorage) RemoveDomainInformation(domainInformationMessages []*models.DomainInformationMessage) error {
	coll := dbWoker.domainInformationCollection()
	transaction := coll.Bulk()
	transaction.Unordered()
	for _, information := range domainInformationMessages {
		transaction.Remove(bson.M{"broker.id": information.Broker.ID, "domain.name": information.RealWorldDomain.Name })
	}
	_, err := transaction.Run()
	return err
}

func (dbWoker *DomainInformationStorage) FindAllDomainInformation() ([]*models.DomainInformationMessage, error) {

	coll := dbWoker.domainInformationCollection()
	var domainInformation []*models.DomainInformationMessage
	var error error

	if error := coll.Find(nil).All(&domainInformation); error != nil {
		fmt.Println(error)
	}
	return domainInformation, error
}

func (dbWoker *DomainInformationStorage) FindDomainInformationByDomainName(domainName string) ([]*models.DomainInformationMessage, error) {
	var domainInformation []*models.DomainInformationMessage
	var error error
	coll := dbWoker.domainInformationCollection()
	regex := bson.M{"$regex":bson.RegEx{"^" + domainName, "m"}}
	fmt.Println(regex)
	if error = coll.Find(bson.M{"domain.name": regex }).All(&domainInformation); error != nil {
		fmt.Println(error)
	}
	return domainInformation, error
}

func (dbWoker *DomainInformationStorage) FindDomainInformationForRequest(informationRequest *models.DomainInformationRequest) ([]*models.DomainInformationMessage, error) {
	var domainInformation []*models.DomainInformationMessage
	var error error
	coll := dbWoker.domainInformationCollection()
	regex := bson.M{"$regex":bson.RegEx{"^" + informationRequest.Domain(), "m"}}
	fmt.Println(regex)
	if error = coll.Find(bson.M{"domain.name": regex, "broker.geolocation.country" : informationRequest.Location().Country }).All(&domainInformation); error != nil {
		fmt.Println(error)
	}
	return domainInformation, error
}

func (dbWorker *DomainInformationStorage) NumberOfBrokersForDomain(domain *models.RealWorldDomain) int {
	count := 0

	coll := dbWorker.domainInformationCollection()
	count, error := coll.Find(bson.M{"domain.name": domain.Name}).Count();
	if error != nil {
		fmt.Println(error)
	}
	return count
}

func (dbWoker *DomainInformationStorage) FindAllBrokers() ([]*models.Broker, error) {
	coll := dbWoker.brokerCollection()
	var domainInformation []*models.DomainInformationMessage
	var error error

	if error := coll.Find(nil).Select(bson.M{"broker":1}).All(&domainInformation); error != nil {
		fmt.Printf("Query Error: %s", error.Error())
	}
	brokers := make([]*models.Broker, 0, len(domainInformation))
	for _, information := range domainInformation {
		brokers = append(brokers, information.Broker)
	}
	return brokers, error
}

func (dbWoker *DomainInformationStorage) StoreDomain(domain *models.RealWorldDomain) error {

	coll := dbWoker.domainCollection()
	index := mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	_ = coll.EnsureIndex(index)
	error := coll.Insert(domain)

	if error != nil && !mgo.IsDup(error) {
		return error
	}
	return nil
}

func (dbWoker *DomainInformationStorage) RemoveDomain(domain *models.RealWorldDomain) error {

	coll := dbWoker.domainCollection()
	error := coll.Remove(domain)

	if error != nil && !mgo.IsDup(error) {
		return error
	}
	return nil
}

func (dbWorker*DomainInformationStorage) FindAllDomains() ([]*models.RealWorldDomain, error) {
	var domains []*models.RealWorldDomain
	coll := dbWorker.domainCollection()
	err := coll.Find(nil).All(&domains)
	return domains, err
}