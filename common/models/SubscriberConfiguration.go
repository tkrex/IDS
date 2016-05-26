package models


type MqttClientConfiguration struct {
	brokerAddress string
	topic 		string
	clientID 	string
}

type MqttProtocol string

const (
	TCP MqttProtocol = "tcp"
	Websocket MqttProtocol ="ws"
)

func NewMqttClientConfiguration(brokerAddress string,port string, protocol MqttProtocol, topic string, clientID string) *MqttClientConfiguration {
	clientConfig := new(MqttClientConfiguration)
	clientConfig.brokerAddress = string(protocol)+"://"+brokerAddress+":"+port
	clientConfig.topic = topic
	clientConfig.clientID = clientID
	return clientConfig
}

func (config *MqttClientConfiguration) BrokerAddress() string {
	return config.brokerAddress
}

func (config *MqttClientConfiguration) Topic() string {
	return config.topic
}

func (config *MqttClientConfiguration) ClientID() string {
	return config.clientID
}
