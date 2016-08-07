package persistence

import (
	"fmt"
	"github.com/tkrex/IDS/common/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
)

const (
	Database = "IDSDaemon"
	TopicCollection = "topics"
	BrokerCollection = "broker"
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
		fmt.Println("DATABASE: ", error)
		return nil, error
	}
	return databaseWorker, error
}

func (dbWoker *DomainInformationStorage)Close() {
	dbWoker.session.Close()
}

func openSession() (*mgo.Session, error) {
	host := os.Getenv("MONGODB_URI")
	session, err := mgo.Dial(host)
	return session, err
}

func (dbWorker *DomainInformationStorage) domainCollection() *mgo.Collection {
	return dbWorker.session.DB(Database).C(DomainCollection)
}

func (dbWorker *DomainInformationStorage) topicCollection() *mgo.Collection {
	return dbWorker.session.DB(Database).C(TopicCollection)
}

func (dbWorker *DomainInformationStorage) brokerCollection() *mgo.Collection {
	return dbWorker.session.DB(Database).C(BrokerCollection)
}

func (dbWorker*DomainInformationStorage) FindDomainInformationByDomainName(domainName string) *models.DomainInformationMessage {

	domainInformation := new(models.DomainInformationMessage)
	topics, topicsError := dbWorker.FindTopicsByDomain(domainName)
	if topicsError != nil {
		fmt.Println(topicsError)
		return nil
	}
	domainInformation.Topics = topics

	broker, brokerError := dbWorker.FindBroker()
	if brokerError != nil {
		fmt.Println(topicsError)
		return nil
	}

	domainInformation.Broker = broker
	domainInformation.RealWorldDomain = models.NewRealWorldDomain(domainName)
	return domainInformation
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
	index := mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	_ = coll.EnsureIndex(index)
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

func (dbWoker *DomainInformationStorage) StoreTopics(topics []*models.TopicInformation)  error {

	coll := dbWoker.topicCollection()
	index := mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	_ = coll.EnsureIndex(index)

	bulkTransaction := coll.Bulk()

	for _, topic := range topics {
		dbWoker.StoreDomain(topic.Domain)
		bulkTransaction.Remove(bson.M{"name": topic.Name })
		bulkTransaction.Insert(topic)

	}
	_, err := bulkTransaction.Run()
	return  err
}



func (dbWoker *DomainInformationStorage) CountTopics() int {
	coll := dbWoker.topicCollection()
	count, err := coll.Find(nil).Count()
	if err != nil {
		count = 0
	}
	return count
}

func (dbWoker *DomainInformationStorage) FindTopicsByNames(topicNames []string) (map[string]*models.TopicInformation, error) {
	coll := dbWoker.topicCollection()
	existingTopics := make(map[string]*models.TopicInformation)
	for _, name := range topicNames {
		var topic models.TopicInformation
		if err := coll.Find(bson.M{"name": name }).One(&topic); err != nil {
			fmt.Println(err)
			continue
		}
		existingTopics[name] = &topic
	}
	return existingTopics, nil
}

func (dbWoker *DomainInformationStorage) FindTopicsByDomain(domainName string) ([]*models.TopicInformation, error) {
	coll := dbWoker.topicCollection()
	topics := []*models.TopicInformation{}
	err := coll.Find(bson.M{"domain.name": domainName }).All(&topics)
	return topics, err
}

func (dbWoker *DomainInformationStorage) StoreBroker(broker *models.Broker) (error) {
	coll := dbWoker.brokerCollection()

	if err := coll.Insert(broker); err != nil {
		return err
	}
	return nil
}

func (dbWoker *DomainInformationStorage) FindBroker() (*models.Broker, error) {
	coll := dbWoker.brokerCollection()

	var error error
	broker := new(models.Broker)
	if error = coll.Find(nil).One(broker); error != nil {
		fmt.Println(error)
	}
	return broker, error
}



