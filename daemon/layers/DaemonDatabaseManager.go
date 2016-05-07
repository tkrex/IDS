package layers

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
	Database = "IDSDaemon"
	TopicCollection = "topics"
	BrokerCollection = "broker"
	DomainControllerCollection = "domainController"
)



type DaemonDatabaseWorker struct {
	session *mgo.Session
}

func NewDaemonDatabaseWorker() (*DaemonDatabaseWorker, error) {
	databaseWorker := new(DaemonDatabaseWorker)
	var error error
	databaseWorker.session, error = openSession()
	if error != nil {
		return nil , error
	}
	return databaseWorker, error
}

func (dbWoker *DaemonDatabaseWorker)Close() {
	dbWoker.session.Close()
}

func openSession() (*mgo.Session,error) {
	session , err := mgo.Dial(Host)
	return session, err
}


func  (dbWoker *DaemonDatabaseWorker) StoreTopics(topics []*models.Topic) error {

	coll := dbWoker.session.DB(Database).C(TopicCollection)
	index := mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	_ = coll.EnsureIndex(index)


	bulkTransaction := coll.Bulk()

	for _,topic := range topics {
		bulkTransaction.Remove(bson.M{"name": topic.Name })
		bulkTransaction.Insert(topic)

	}
	bulkResult, err := bulkTransaction.Run()
	if err != nil {
		return err
	}
	fmt.Println(bulkResult)
	return nil
}

func  (dbWoker *DaemonDatabaseWorker) StoreTopic(topic *models.Topic) error {

	coll := dbWoker.session.DB(Database).C(TopicCollection)
	index := mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	_ = coll.EnsureIndex(index)
	error := coll.Insert(topic)

	if error != nil {
		return error
	}
	return nil
}

func (dbWoker *DaemonDatabaseWorker)FindAllTopics() ([]*models.Topic,error) {

	coll := dbWoker.session.DB(Database).C(TopicCollection)
	var topics []*models.Topic
	if err := coll.Find(nil).All(&topics); err != nil {
		fmt.Println(err)
	}
	return topics, nil
}

func (dbWoker *DaemonDatabaseWorker) FindTopicsByName(topicNames []string) (map[string]*models.Topic,error) {
	coll := dbWoker.session.DB(Database).C(TopicCollection)
	existingTopics := make(map[string]*models.Topic)
	for _,name := range topicNames {
		var topic models.Topic
		if err := coll.Find(bson.M{"name": name }).One(&topic); err != nil {
			fmt.Println(err)
			continue
		}
		existingTopics[name] = &topic
	}
	return existingTopics, nil
}

func (dbWoker *DaemonDatabaseWorker) StoreBroker(broker *models.Broker) (error) {
	coll := dbWoker.session.DB(Database).C(BrokerCollection)

	if err := coll.Insert(broker); err != nil {
		return err
	}
	return nil
}

func (dbWoker *DaemonDatabaseWorker) FindBroker() (*models.Broker,error) {
	coll := dbWoker.session.DB(Database).C(BrokerCollection)

	var error error
	broker := new(models.Broker)
	if error = coll.Find(nil).One(broker); error != nil {
		fmt.Println(error)
	}
	return broker, error
}


func (dbWoker *DaemonDatabaseWorker) StoreDomainControllers(domainControllers []*models.DomainController) error {
	coll := dbWoker.session.DB(Database).C(DomainControllerCollection)
	bulk := coll.Bulk()
	bulk.Unordered()
	for _, domainController := range domainControllers {
		bulk.Upsert(bson.M{"domain.name":domainController.Domain.Name},bson.M{"$set": domainController})
	}
	_, error := bulk.Run()
	return error
}


