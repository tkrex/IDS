package daemon

import (
	"sync"
	"time"
	"github.com/tkrex/IDS/daemon/layers"
	"github.com/tkrex/IDS/daemon/layers/producer"
	"github.com/tkrex/IDS/daemon/layers/consumer"
	"github.com/tkrex/IDS/daemon/models"
	"github.com/golang/glog"
)

type IDSDaemon struct {
	// Flag to indicate that the session state
	// 0 == Running
	// 1 == Closed
	state int64

	incomingTopicChannel chan models.Topic

	processorStarted sync.WaitGroup
	processorStopped sync.WaitGroup

	// Interval bounds
	publishFrequency            time.Duration
	expectedLast                time.Time


	// Statistics
	totalFrequencyBlocks        uint64
	totalFrequencySleepDuration uint64
	totalFrequencyMisses        uint64

	producer layers.TopicProducer
	consumer layers.TopicConsumer
}


func NewIDSDaemon() *IDSDaemon {
	daemon :=  new(IDSDaemon)
	daemon.incomingTopicChannel = make(chan models.Topic, 100)

	daemon.processorStarted.Add(1)
	daemon.processorStopped.Add(1)
	go daemon.processor()


	daemon.processorStarted.Wait()

	return daemon
}




func (daemon *IDSDaemon) processor() {

	// Talk to db
	daemon.producer = gathering.NewMqttTopicCollector("tcp://127.0.0.1:1883",daemon.incomingTopicChannel)
	daemon.consumer = persistance.NewDataManager(daemon.incomingTopicChannel)


	// Notify constructor
	daemon.processorStarted.Done()
}

func (daemon *IDSDaemon) Run () {
	for consumerStopped, producerStopped  := daemon.consumer.State() == 1, daemon.producer.State() == 1; !(consumerStopped && producerStopped); consumerStopped, producerStopped  = daemon.consumer.State() == 1, daemon.producer.State() == 1 {

	}
	glog.Info("Daemon Stopped")
}
