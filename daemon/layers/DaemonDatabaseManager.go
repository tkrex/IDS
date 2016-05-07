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
)



func openSession() (*mgo.Session,error) {
	session , err := mgo.Dial(Host)
	return session, err
}

func isDatabaseAvailable() bool {
	session, err := openSession()
	defer session.Close()
	if err != nil {
		return false
	}
	return true
}

func  StoreTopics(topics []*models.Topic) error {

	session, err := openSession()
	if err != nil {
		return err
	}

	defer session.Close()
	fmt.Printf("Connected to %v\n", session.LiveServers())

	coll := session.DB(Database).C(TopicCollection)
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

func  StoreTopic(topic *models.Topic) error {

	session,err := openSession()
	if err != nil {
		return err
	}

	defer session.Close()
	fmt.Printf("Connected to %v\n", session.LiveServers())

	coll := session.DB(Database).C(TopicCollection)
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

func FindAllTopics() ([]*models.Topic,error) {
	session,err := openSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	coll := session.DB(Database).C(TopicCollection)
	var topics []*models.Topic
	if err := coll.Find(nil).All(&topics); err != nil {
		fmt.Println(err)
	}
	return topics, nil
}

func FindTopicsByName(topicNames []string) (map[string]*models.Topic,error) {
	session,err := openSession()
	if err != nil {
		return nil, err
	}

	defer session.Close()
	coll := session.DB(Database).C(TopicCollection)
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

func StoreBroker(broker *models.Broker) (error) {
	session,err := openSession()
	if err != nil {
		return err
	}

	defer session.Close()

	fmt.Printf("Connected to %v\n", session.LiveServers())

	coll := session.DB(Database).C(BrokerCollection)

	if err := coll.Insert(broker); err != nil {
		return err
	}
	return nil
}

func FindBroker() (*models.Broker,error) {


	session,err := openSession()
	if err != nil {
		return nil,err
	}

	defer session.Close()
	coll := session.DB(Database).C(BrokerCollection)

	var error error
	broker := new(models.Broker)
	if error = coll.Find(nil).One(broker); error != nil {
		fmt.Println(error)
	}
	return broker, error
}


