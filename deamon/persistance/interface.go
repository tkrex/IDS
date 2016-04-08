package persistance
import (
	"persistance"
)
type TopicPersistanceDelagte interface
	
	Store(topic Topic)

	FindTopicByName(name string) Topic

	FindAllTopics() []Topic
	

