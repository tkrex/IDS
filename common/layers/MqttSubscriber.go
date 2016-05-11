package common
import (
  "fmt"
  //import the Paho Go MQTT library
  MQTT "github.com/eclipse/paho.mqtt.golang"
  "os"
  "sync"
  "sync/atomic"
  "github.com/tkrex/IDS/common/models"
)

type  MqttSubscriber struct {

  state                 int64
  incomingTopicsChannel chan *models.RawTopicMessage
  client                MQTT.Client

  producerStarted       sync.WaitGroup
  producerStopped       sync.WaitGroup

  topicCounter          int
  config                *models.MqttClientConfiguration
  isDaemon              bool
}




func NewMqttSubscriber(subscriberConfig *models.MqttClientConfiguration, isDaemon bool) *MqttSubscriber {
  subscriber := new(MqttSubscriber)
  subscriber.incomingTopicsChannel = make(chan *models.RawTopicMessage,100)
  subscriber.config = subscriberConfig
  subscriber.isDaemon = isDaemon
  subscriber.producerStarted.Add(1)
  subscriber.producerStopped.Add(1)
  opts := MQTT.NewClientOptions().AddBroker(subscriber.config.BrokerAddress())
  opts.SetClientID(subscriberConfig.ClientID())
      //define a function for the default message handler
  var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
      subscriber.onReceiveMessage(msg)
    }
  opts.SetDefaultPublishHandler(f)

  subscriber.client = MQTT.NewClient(opts)

  go subscriber.run()
  subscriber.producerStarted.Wait()
  return subscriber
}

func (subscriber *MqttSubscriber) IncomingTopicsChannel() chan *models.RawTopicMessage {
  return subscriber.incomingTopicsChannel
}


func (subscriber *MqttSubscriber) State() int64 {
  return atomic.LoadInt64(&subscriber.state)
}

func (subscriber *MqttSubscriber) run() {
  
  if token := subscriber.client.Connect(); token.Wait() && token.Error() != nil {
    panic(token.Error())
  }

  fmt.Println("Connected to Mqtt Cient")
  //subscribe to the topic /go-mqtt/sample and request messages to be delivered
  //at a maximum qos of zero, wait for the receipt to confirm the subscription
  if token := subscriber.client.Subscribe(subscriber.config.Topic(), 0, nil); token.Wait() && token.Error() != nil {
    fmt.Println(token.Error())
    subscriber.Close()
  }
  subscriber.producerStarted.Done()
}

func (subscriber *MqttSubscriber) stopCollectingTopics() {
  defer subscriber.producerStopped.Done()

  fmt.Println("Unsubscribing")
  if token :=  subscriber.client.Unsubscribe(subscriber.config.Topic()); token.Wait() && token.Error() != nil {
    fmt.Println(token.Error())
    os.Exit(1)
  }
  fmt.Println("Unsubscribed")

  subscriber.client.Disconnect(10)
  fmt.Println("Disconnected")
}


func (subscriber *MqttSubscriber) onReceiveMessage(msg MQTT.Message) {
  rawTopic := models.NewRawTopicMessage(msg.Topic(),msg.Payload())

  //DEBUG in case of one local broker for daemon and domainController
  if subscriber.isDaemon && rawTopic.Name == "domainController" {
    return
  }
  if closed := subscriber.State() == 1; !closed {
    fmt.Println(rawTopic.Name)
    subscriber.incomingTopicsChannel <- rawTopic
  }

}

func (subscriber *MqttSubscriber) Close() {
  fmt.Println("Closing Subscriber")
  atomic.StoreInt64(&subscriber.state,1)
  close(subscriber.incomingTopicsChannel)
  subscriber.stopCollectingTopics()
  fmt.Println("Change Producer State")
  subscriber.producerStopped.Wait()
}