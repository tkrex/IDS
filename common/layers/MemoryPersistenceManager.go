package common

import "github.com/tkrex/IDS/common/models"

type MemoryPersistenceManager struct {
	topics map[string]*models.Topic
}

func NewMemoryPersistenceManager() *MemoryPersistenceManager {
	manager := new(MemoryPersistenceManager)
	manager.topics = make(map[string]*models.Topic)
	return manager
}


func (manager *MemoryPersistenceManager) Topics() []*models.Topic {
	topicArray  := make([]*models.Topic,len(manager.topics))

	index := 0
	for _,topic := range manager.topics {
		topicArray[index] = topic
		index++
	}
	return topicArray
}

func (manager *MemoryPersistenceManager) TopicWithName(name string) (*models.Topic,bool) {
	topic, ok := manager.topics[name]
	return topic,ok
}

func (manager *MemoryPersistenceManager) StoreTopic(topic *models.Topic) {
	manager.topics[topic.Name] = topic
}