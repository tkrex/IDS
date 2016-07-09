package persistence

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"github.com/tkrex/IDS/common/models"
	"gopkg.in/mgo.v2/bson"
	"os"
)

const (
	Database = "IDSDomainController"
	BrokerCollection = "brokers"
	DomainInformationCollection = "domainInformation"
	DomainCollection = "domains"
)

type DomainControllerDatabaseWorker struct {
	session *mgo.Session
}

func NewDomainControllerDatabaseWorker() (*DomainControllerDatabaseWorker, error) {
	databaseWorker := new(DomainControllerDatabaseWorker)
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

func (dbWoker *DomainControllerDatabaseWorker)Close() {
	dbWoker.session.Close()
}

func (dbWorker *DomainControllerDatabaseWorker) domainInformationCollection() *mgo.Collection {
	return dbWorker.session.DB(Database).C(DomainInformationCollection)
}

func (dbWorker *DomainControllerDatabaseWorker) domainCollection() *mgo.Collection {
	return dbWorker.session.DB(Database).C(DomainCollection)
}

func (dbWorker *DomainControllerDatabaseWorker) brokerCollection() *mgo.Collection {
	return dbWorker.session.DB(Database).C(BrokerCollection)
}

func (dbWoker *DomainControllerDatabaseWorker) StoreBroker(broker *models.Broker) (error) {
	coll := dbWoker.brokerCollection()

	if err := coll.Insert(broker); err != nil {
		return err
	}
	return nil
}

func (dbWoker *DomainControllerDatabaseWorker) FindBroker() (*models.Broker, error) {
	coll := dbWoker.brokerCollection()

	var error error
	broker := new(models.Broker)
	if error = coll.Find(nil).One(broker); error != nil {
		fmt.Println(error)
	}
	return broker, error
}

func (dbWoker *DomainControllerDatabaseWorker) StoreDomainInformationMessages(domainInformationMessages []*models.DomainInformationMessage) error {
	domainInformationCollection := dbWoker.domainInformationCollection()
	index := mgo.Index{
		Key:        []string{"broker.id", "domain.name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	_ = domainInformationCollection.EnsureIndex(index)

	dbWoker.RemoveDomainInformation(domainInformationMessages)
	insertDomainInformationBulk := domainInformationCollection.Bulk()

	for _, information := range domainInformationMessages {
		insertDomainInformationBulk.Insert(information)
		dbWoker.StoreDomain(information.RealWorldDomain)

	}
	_, err := insertDomainInformationBulk.Run()

	return err
}

func (dbWoker *DomainControllerDatabaseWorker) StoreDomainInformationMessage(domainInformationMessage *models.DomainInformationMessage) error {
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

func (dbWoker *DomainControllerDatabaseWorker) RemoveDomainInformation(domainInformationMessages []*models.DomainInformationMessage) error {
	coll := dbWoker.domainInformationCollection()
	transaction := coll.Bulk()
	transaction.Unordered()
	for _, information := range domainInformationMessages {
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
	regex := bson.M{"$regex":bson.RegEx{"^" + domainName, "m"}}
	fmt.Println(regex)
	if error = coll.Find(bson.M{"domain.name": regex }).All(&domainInformation); error != nil {
		fmt.Println(error)
	}
	return domainInformation, error
}

func (dbWoker *DomainControllerDatabaseWorker) FindDomainInformationForRequest(informationRequest *models.DomainInformationRequest) ([]*models.DomainInformationMessage, error) {
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

func (dbWoker *DomainControllerDatabaseWorker) FindDomainInformationForBroker(informationRequest *models.DomainInformationRequest, brokerId string) ([]*models.DomainInformationMessage, error) {
	var domainInformation []*models.DomainInformationMessage
	var error error
	regex := bson.M{"$regex":bson.RegEx{"^" + informationRequest.Domain(), "m"}}

	coll := dbWoker.domainInformationCollection()
	if error = coll.Find(bson.M{"domain.name": regex, "broker.id": brokerId, "broker.geolocation.country" : informationRequest.Location().Country }).All(&domainInformation); error != nil {
		fmt.Println(error)
	}
	return domainInformation, error
}

func (dbWoker *DomainControllerDatabaseWorker) FindBrokersForDomain(domainName string) ([]*models.Broker, error) {
	coll := dbWoker.brokerCollection()
	var domainInformation []*models.DomainInformationMessage
	var error error
	regex := bson.M{"$regex":bson.RegEx{"^" + domainName, "m"}}
	if error := coll.Find(bson.M{"domain.name": regex }).Select(bson.M{"broker":1}).All(&domainInformation); error != nil {
		fmt.Printf("Query Error: %s", error.Error())
	}
	brokers := make([]*models.Broker, 0, len(domainInformation))
	for _, information := range domainInformation {
		brokers = append(brokers, information.Broker)
	}
	return brokers, error
}

func (dbWoker *DomainControllerDatabaseWorker) FindBrokersForInformationRequest(informationRequest *models.DomainInformationRequest) ([]*models.Broker, error) {
	coll := dbWoker.brokerCollection()
	var domainInformation []*models.DomainInformationMessage
	var error error
	regex := bson.M{"$regex":bson.RegEx{"^" + informationRequest.Domain(), "m"}}
	if error := coll.Find(bson.M{"domain.name": regex, "broker.geolocation.country":informationRequest.Location().Country}).Select(bson.M{"broker":1}).All(&domainInformation); error != nil {
		fmt.Printf("Query Error: %s", error.Error())
		return nil,error
	}
	brokers := make([]*models.Broker, 0, len(domainInformation))
	for _, information := range domainInformation {
		brokers = append(brokers, information.Broker)
	}
	return brokers, error
}

func (dbWoker *DomainControllerDatabaseWorker) FindAllBrokers() ([]*models.Broker, error) {
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

func (dbWoker *DomainControllerDatabaseWorker) StoreDomain(domain *models.RealWorldDomain) error {

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

func (dbWoker *DomainControllerDatabaseWorker) RemoveDomain(domain *models.RealWorldDomain) error {

	coll := dbWoker.domainCollection()
	error := coll.Remove(domain)

	if error != nil && !mgo.IsDup(error) {
		return error
	}
	return nil
}

func (dbWorker*DomainControllerDatabaseWorker) FindAllDomains() ([]*models.RealWorldDomain, error) {
	var domains []*models.RealWorldDomain
	coll := dbWorker.domainCollection()
	err := coll.Find(nil).All(&domains)
	return domains, err
}