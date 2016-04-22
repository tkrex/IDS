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
  incomingTopicsChannel chan *models.Topic
  client                MQTT.Client

  consumer              InformationProcessor
  producerStarted       sync.WaitGroup
  producerStopped       sync.WaitGroup

  topicCounter          int
  desiredTopic          string
}


func (collector *MqttSubscriber) InformationChannel() chan *models.Topic {
  return collector.incomingTopicsChannel
}


func NewMqttSubscriber(brokerAddress string, desiredTopic string) *MqttSubscriber {
  subscriber := new(MqttSubscriber)
  subscriber.incomingTopicsChannel = make(chan *models.Topic,100)
  subscriber.desiredTopic = desiredTopic
  subscriber.producerStarted.Add(1)
  subscriber.producerStopped.Add(1)
  opts := MQTT.NewClientOptions().AddBroker(brokerAddress)
  opts.SetClientID("subscriber")
      //define a function for the default message handler
  var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
      subscriber.OnReceiveMessage(msg)
    }
  opts.SetDefaultPublishHandler(f)

  subscriber.client = MQTT.NewClient(opts)

  go subscriber.processor()
  subscriber.producerStarted.Wait()
  return subscriber
}

func (subscriber* MqttSubscriber) createConsumer() {
  subscriber.consumer = NewTopicProcessor(subscriber)
}

func (subscriber * MqttSubscriber) CloseConsumer() {
  subscriber.consumer.Close()
}

func (subscriber *MqttSubscriber) State() int64 {
  return atomic.LoadInt64(&subscriber.state)
}

func (subscriber *MqttSubscriber) processor() {
  
  if token := subscriber.client.Connect(); token.Wait() && token.Error() != nil {
    panic(token.Error())
  }

  //Create Consumer and wait for it being ready to consume topics
  subscriber.createConsumer()
  fmt.Println("Connected to Mqtt Cient")
  //subscribe to the topic /go-mqtt/sample and request messages to be delivered
  //at a maximum qos of zero, wait for the receipt to confirm the subscription
  if token := subscriber.client.Subscribe(subscriber.desiredTopic, 0, nil); token.Wait() && token.Error() != nil {
    fmt.Println(token.Error())
    os.Exit(1)
  }
  subscriber.producerStarted.Done()
}

func (subscriber *MqttSubscriber) StopCollectingTopics() {
  defer subscriber.producerStopped.Done()

  fmt.Println("Unsubscribing")
  if token :=  subscriber.client.Unsubscribe(subscriber.desiredTopic); token.Wait() && token.Error() != nil {
    fmt.Println(token.Error())
    os.Exit(1)
  }
  fmt.Println("Unsubscribed")

  subscriber.client.Disconnect(10)
  fmt.Println("Disconnected")
}


func (subscriber *MqttSubscriber) OnReceiveMessage(msg MQTT.Message) {
  topic := models.NewTopic(1,msg.Topic(), msg.Payload())
  if closed := subscriber.State() == 1; !closed {
    fmt.Println(topic.Name)
    subscriber.incomingTopicsChannel <- topic
  }

}

func (subscriber *MqttSubscriber) Close() {
  fmt.Println("Closing Collector")
  close(subscriber.incomingTopicsChannel)
  subscriber.StopCollectingTopics()
  subscriber.consumer.Close()
  fmt.Println("Change Producer State")
  atomic.StoreInt64(&subscriber.state,1)
  subscriber.producerStopped.Wait()
}