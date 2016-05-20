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
	DomainCollection = "domains"
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

func (dbWorker *DaemonDatabaseWorker) domainCollection() *mgo.Collection {
	return dbWorker.session.DB(Database).C(DomainCollection)
}

func (dbWorker *DaemonDatabaseWorker) topicCollection() *mgo.Collection {
	return dbWorker.session.DB(Database).C(TopicCollection)
}

func (dbWorker *DaemonDatabaseWorker) brokerCollection() *mgo.Collection {
	return dbWorker.session.DB(Database).C(BrokerCollection)
}







func (dbWorker* DaemonDatabaseWorker) FindDomainInformationByDomainName(domainName string) (*models.DomainInformationMessage){

	 domainInformation := new(models.DomainInformationMessage)
	topics ,topicsError :=dbWorker.FindTopicsByDomain(domainName)
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

func (dbWorker* DaemonDatabaseWorker) FindAllDomainInformation() ([]*models.DomainInformationMessage,error){

	domains, err := dbWorker.FindAllDomains()
	if err != nil {
		return nil, err
	}
	var domainInformationMessages []*models.DomainInformationMessage
	for _,domain := range domains {
		domainInformation := dbWorker.FindDomainInformationByDomainName(domain.Name)
		if domainInformation != nil {
			domainInformationMessages = append(domainInformationMessages, domainInformation)
		}
	}
	return domainInformationMessages, nil
}



func  (dbWoker *DaemonDatabaseWorker) StoreDomain(domain *models.RealWorldDomain) error {

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


func  (dbWoker *DaemonDatabaseWorker) RemoveDomain(domain *models.RealWorldDomain) error {

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

func (dbWorker* DaemonDatabaseWorker) FindAllDomains() ([]*models.RealWorldDomain,error){
	var domains []*models.RealWorldDomain
	coll := dbWorker.domainCollection()
	err := coll.Find(nil).All(&domains)
	return domains, err
}



func  (dbWoker *DaemonDatabaseWorker) StoreTopics(topics []*models.Topic) (*mgo.BulkResult,error) {

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

	for _,topic := range topics {
		dbWoker.StoreDomain(topic.Domain)
		bulkTransaction.Remove(bson.M{"name": topic.Name })
		bulkTransaction.Insert(topic)

	}
	bulkResult, err := bulkTransaction.Run()
	return bulkResult,err
}


func  (dbWoker *DaemonDatabaseWorker) StoreTopic(topic *models.Topic) error {

	coll := dbWoker.topicCollection()
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

	coll := dbWoker.topicCollection()
	var topics []*models.Topic
	if err := coll.Find(nil).All(&topics); err != nil {
		fmt.Println(err)
	}
	return topics, nil
}

func (dbWoker *DaemonDatabaseWorker) FindTopicsByName(topicNames []string) (map[string]*models.Topic,error) {
	coll := dbWoker.topicCollection()
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

func (dbWoker *DaemonDatabaseWorker) FindTopicsByDomain(domainName string) ([]*models.Topic,error) {
	coll := dbWoker.topicCollection()
	topics := []*models.Topic{}
	 err := coll.Find(bson.M{"domain.name": domainName }).All(&topics)
	return topics, err
}

func (dbWoker *DaemonDatabaseWorker) StoreBroker(broker *models.Broker) (error) {
	coll := dbWoker.brokerCollection()

	if err := coll.Insert(broker); err != nil {
		return err
	}
	return nil
}

func (dbWoker *DaemonDatabaseWorker) FindBroker() (*models.Broker,error) {
	coll := dbWoker.brokerCollection()

	var error error
	broker := new(models.Broker)
	if error = coll.Find(nil).One(broker); error != nil {
		fmt.Println(error)
	}
	return broker, error
}



