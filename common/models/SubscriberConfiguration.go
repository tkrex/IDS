package models

import "net/url"

type MqttClientConfiguration struct {
	brokerAddress url.URL
	topic 		string
	clientID 	string
}

type MqttProtocol string

const (
	TCP MqttProtocol = "tcp"
	Websocket MqttProtocol ="ws"
)

func NewMqttClientConfiguration(brokerAddress url.URL, topic string, clientID string) *MqttClientConfiguration {
	clientConfig := new(MqttClientConfiguration)
	clientConfig.brokerAddress = brokerAddress
	clientConfig.topic = topic
	clientConfig.clientID = clientID
	return clientConfig
}

func (config *MqttClientConfiguration) BrokerAddress() string {
	return config.brokerAddress.String()
}

func (config *MqttClientConfiguration) Topic() string {
	return config.topic
}

func (config *MqttClientConfiguration) ClientID() string {
	return config.clientID
}
