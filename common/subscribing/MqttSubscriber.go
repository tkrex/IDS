package subscribing
import (
  "fmt"
  //import the Paho Go MQTT library
  MQTT "github.com/eclipse/paho.mqtt.golang"
  "os"
  "sync"
  "sync/atomic"
  "github.com/tkrex/IDS/common/models"
)
//Provides the functionality of subscribing to topics of a MQTT Broker
type  MqttSubscriber struct {

  state                 int64
  incomingTopicsChannel chan *models.RawTopicMessage
  client                MQTT.Client

  topic string

  producerStarted       sync.WaitGroup
  producerStopped       sync.WaitGroup
  config                *models.MqttClientConfiguration
}




func NewMqttSubscriber(subscriberConfig *models.MqttClientConfiguration, topic string) *MqttSubscriber {
  subscriber := new(MqttSubscriber)
  subscriber.incomingTopicsChannel = make(chan *models.RawTopicMessage,100)
  subscriber.config = subscriberConfig
  subscriber.topic = topic
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
//Returns a channel to which the MqttSubscriber forwards incoming MQTT messages
func (subscriber *MqttSubscriber) IncomingTopicsChannel() chan *models.RawTopicMessage {
  return subscriber.incomingTopicsChannel
}

//Returns the current state of the Subscriber
func (subscriber *MqttSubscriber) State() int64 {
  return atomic.LoadInt64(&subscriber.state)
}

//Connects to the MQTT Broker specified in the MqttClientConfig and subscribes to the topic specified when initializing the subscriber
func (subscriber *MqttSubscriber) run() {
  
  if token := subscriber.client.Connect(); token.Wait() && token.Error() != nil {
    panic(token.Error())
  }

  fmt.Println("Connected to Mqtt Cient")
  //subscribe to the topic /go-mqtt/sample and request messages to be delivered
  //at a maximum qos of zero, wait for the receipt to confirm the subscription
  if token := subscriber.client.Subscribe(subscriber.topic, 0, nil); token.Wait() && token.Error() != nil {
    fmt.Println(token.Error())
    subscriber.Close()
  }
  subscriber.producerStarted.Done()
}


//Unsubscribes from the topic and disconnects from the MQTT Broker
func (subscriber *MqttSubscriber) stopCollectingTopics() {
  defer subscriber.producerStopped.Done()

  fmt.Println("Unsubscribing")
  if token :=  subscriber.client.Unsubscribe(subscriber.topic); token.Wait() && token.Error() != nil {
    fmt.Println(token.Error())
    os.Exit(1)
  }
  fmt.Println("Unsubscribed")

  subscriber.client.Disconnect(10)
  fmt.Println("Disconnected")
}

//Callback for receiving MQTT messages
func (subscriber *MqttSubscriber) onReceiveMessage(msg MQTT.Message) {
  rawTopic := models.NewRawTopicMessage(msg.Topic(),msg.Payload())

  if closed := subscriber.State() == 1; !closed {
    fmt.Println(rawTopic.Topic)
    subscriber.incomingTopicsChannel <- rawTopic
  }

}

//Closes incomingTopicsChannel and disconnects from MQTT Broker
func (subscriber *MqttSubscriber) Close() {
  fmt.Println("Closing Subscriber")
  atomic.StoreInt64(&subscriber.state,1)
  close(subscriber.incomingTopicsChannel)
  subscriber.stopCollectingTopics()
  fmt.Println("Change Producer State")
  subscriber.producerStopped.Wait()
}