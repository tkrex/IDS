package publishing

import (
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/tkrex/IDS/common/models"
	"fmt"
)
//Implements the functionality to publish data to a MQTT Broker
type MqttPublisher struct {
	client                mqtt.Client
	config *models.MqttClientConfiguration
	retained bool
}

func NewMqttPublisher(config *models.MqttClientConfiguration, retained bool) *MqttPublisher {
	publisher := new(MqttPublisher)
	publisher.config = config
	publisher.retained = retained
	publisher.connect()
	fmt.Println("Publisher started")
	return publisher
}

//Connects with the MQTT Broker, specified in the MqttClientConfiguration
func (publisher *MqttPublisher) connect() {
	opts := mqtt.NewClientOptions().AddBroker(publisher.config.BrokerAddress())
	opts.SetClientID(publisher.config.ClientID())
	opts.WillRetained = publisher.retained
	publisher.client = mqtt.NewClient(opts)

	if token := publisher.client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

//Publishers data to a topic
func (publisher *MqttPublisher) Publish(data []byte, topic string) error  {
	if token := publisher.client.Publish(topic, 2, false, data); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	fmt.Println("Message Published")
	return nil
}

//Disconnects from the MQTT Broker
func (publisher *MqttPublisher) Close() {
	publisher.client.Disconnect(10)
	fmt.Println("Publisher Disconnected")
}