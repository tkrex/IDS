package models


type MqttClientConfiguration struct {
	brokerAddress string
	topic 		string
	clientID 	string
}

func NewMqttClientConfiguration(brokerAddress string, topic string, clientID string) *MqttClientConfiguration {
	clientConfig := new(MqttClientConfiguration)
	clientConfig.brokerAddress = "tcp://"+brokerAddress+":1883"
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
