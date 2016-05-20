package common

import (
	"github.com/tkrex/IDS/common/models"
)

type MemoryTopicPersistenceManager struct {
	topics map[string] models.Topic
}

func NewMemoryTopicPersistenceManager() *MemoryTopicPersistenceManager {
	manager := new(MemoryTopicPersistenceManager)
	manager.topics = make(map[string]models.Topic)
	return manager
}


func (manager *MemoryTopicPersistenceManager) Topics() []models.Topic {
	topicArray  := make([]models.Topic,len(manager.topics))

	index := 0
	for _,topic := range manager.topics {
		topicArray[index] = topic
		index++
	}
	return topicArray
}

func (manager *MemoryTopicPersistenceManager) TopicWithName(name string) (models.Topic,bool) {
	topic, ok := manager.topics[name]
	return topic,ok
}

func (manager *MemoryTopicPersistenceManager) StoreTopic(topic models.Topic) {
	manager.topics[topic.Name] = topic
}

func (manager *MemoryTopicPersistenceManager) NumberOfTopics() int  {
	return len(manager.topics)
}