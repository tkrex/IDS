package layers
import "github.com/tkrex/IDS/daemon/models"

type TopicConsumer interface {
	Store(topic models.Topic)
	State() int64
}
 
type TopicProducer interface {
	State() int64
}