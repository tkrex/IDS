package layers

import (
	"github.com/tkrex/IDS/common/models"
	"encoding/json"
	"net/http"
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"github.com/tkrex/IDS/daemon"
)

type BrokerRegistrationWorker struct {
	broker *models.Broker
}


func NewBrokerRegistrationWorker() *BrokerRegistrationWorker {
	worker := new(BrokerRegistrationWorker)
	worker.broker = new(models.Broker)
	return worker
}

func (worker *BrokerRegistrationWorker) GatherBrokerInformation() {
	//TODO: Get IP address from Docker ENV
	worker.broker.IP = "66.220.158.68"
	worker.findBrokerDomainName()
	worker.findBrokerGeolocation()
	fmt.Print(worker.broker)
}


func (worker *BrokerRegistrationWorker) findBrokerDomainName() {
	name, err := net.LookupAddr(worker.broker.IP)
	if err != nil {
		worker.broker.InternetDomain = ""
		return
	}
	worker.broker.InternetDomain = name[0]
}

func (worker *BrokerRegistrationWorker) findBrokerGeolocation() {
	geolocationFetcher := common.NewGeoLocationFetcher("192.168.99.100")
	location , err := geolocationFetcher.SendGeoLocationRequest(worker.broker.IP)
	if err != nil {
		worker.broker.Geolocation = new(models.Geolocation)
		return
	}
	worker.broker.Geolocation = location
}

func (worker *BrokerRegistrationWorker) registerBroker() {
	//TODO: get own Broker information
	broker := models.NewBroker("8.8.8.8", "krex.com")

	jsonString, _ := json.Marshal(&broker)

	req, err := http.NewRequest("POST", "http://localhost:8080/rest/brokers", bytes.NewBuffer(jsonString))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		fmt.Println(body)
		return
	}

	reponseBroker := new(models.Broker)

	err = json.Unmarshal(body,reponseBroker)
	if err != nil {
		fmt.Println("REGISTER BROKER: Unkown response format")
		return
	}

	//TODO: store new broker ID



}