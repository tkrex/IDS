package persistance

import (
	"github.com/tkrex/IDS/daemon/models"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
    "fmt"
)

type TopicDataManager struct {
	database *gorm.DB
}

func NewTopicDataManager() *TopicDataManager {
	manager := new(TopicDataManager)
	manager.CreateDatabase()
	manager.CreateTopicsTable()
	return manager
}

func (manager *TopicDataManager) CreateDatabase() {
	 db, err := gorm.Open("sqlite3", "./topics.db")
	 if err != nil {
	 	 panic("failed to connect database")
	 }
	 manager.database = db
}

func (manager *TopicDataManager) CreateTopicsTable() {
	if !manager.database.HasTable(&models.Topic{}) {
		manager.database.CreateTable(&models.Topic{})
	}
}

func (manager *TopicDataManager) Store(topic *models.Topic) {
	var count int
	manager.database.Model(&models.Topic{}).Where("name = ?", topic.Name).Count(&count) 
	fmt.Println(count)
	if count == 0 {
		fmt.Println("new Topic")
		manager.database.Create(&topic)
	} else {
		manager.UpdateTopic(topic)
		fmt.Println("Update Topic")
	}
}

func (manager *TopicDataManager) UpdateTopic(topic *models.Topic) {
	var oldTopic models.Topic
	manager.database.Model(&models.Topic{}).Where("name = ?", topic.Name).First(&oldTopic)
	fmt.Println(oldTopic.LastUpdateTimeStamp)
	fmt.Println(topic.LastUpdateTimeStamp)

	updateInterval := topic.LastUpdateTimeStamp.Sub(oldTopic.LastUpdateTimeStamp).Seconds()
	fmt.Println(updateInterval)
	manager.database.Model(&topic).UpdateColumn("updateInterval", updateInterval)
	manager.database.Model(&topic).UpdateColumn("LastUpdateTimeStamp", topic.LastUpdateTimeStamp)
	manager.database.Model(&topic).UpdateColumn("LastPayload", topic.LastPayload)
}