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
	Collection = "topics"
)



func OpenSession() *mgo.Session {
	//session, err := mgo.DialWithInfo(&mgo.DialInfo{
	//	Addrs:    []string{Host},
	//	Username: Username,
	//	Password: Password,
	//	Database: Database,
	//	DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
	//		return tls.Dial("tcp", addr.String(), &tls.Config{})
	//	},
	//})

	session , err := mgo.Dial(Host)

	if err != nil {
		panic(err)
	}

	return session
}
func  StoreTopics(topics []*models.Topic) {

	session := OpenSession()

	defer session.Close()
	fmt.Printf("Connected to %v\n", session.LiveServers())

	coll := session.DB(Database).C(Collection)
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
		fmt.Println(err)
	}
	fmt.Println(bulkResult)
}

func  StoreTopic(topic *models.Topic) {

	session := OpenSession()

	defer session.Close()
	fmt.Printf("Connected to %v\n", session.LiveServers())

	coll := session.DB(Database).C(Collection)
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
		fmt.Println(error)
	}
}

func FindAllTopics() []*models.Topic {
	session := OpenSession()
	defer session.Close()
	coll := session.DB(Database).C(Collection)
	var topics []*models.Topic
	if err := coll.Find(nil).All(&topics); err != nil {
		fmt.Println(err)
	}
	return topics
}

func FindTopicsByName(topicNames []string) map[string]*models.Topic {
	session := OpenSession()
	defer session.Close()
	coll := session.DB(Database).C(Collection)
	existingTopics := make(map[string]*models.Topic)
	for _,name := range topicNames {
		var topic models.Topic
		if err := coll.Find(bson.M{"name": name }).One(&topic); err != nil {
			fmt.Println(err)
			continue
		}
		existingTopics[name] = &topic
	}
	return existingTopics
}

