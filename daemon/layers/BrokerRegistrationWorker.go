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
	"time"
	"sync"
)

const RegisterInterval  = time.Second * 10

type BrokerRegistrationWorker struct {
	broker *models.Broker
	registerTicker *time.Ticker
	workerStarted sync.WaitGroup
	workerStopped sync.WaitGroup
}

func NewBrokerRegistrationWorker() *BrokerRegistrationWorker {
	worker := new(BrokerRegistrationWorker)
	worker.broker = new(models.Broker)
	worker.workerStarted.Add(1)
	worker.workerStopped.Add(2)
	go worker.registerBroker()
	worker.workerStarted.Wait()
	return worker
}

func (worker *BrokerRegistrationWorker) registerBroker() {
	if !isDatabaseAvailable() {
		fmt.Println("Database not reachable")
		return
	}
	worker.workerStarted.Done()

	if isBrokerRegistered := worker.isBrokerRegistered(); isBrokerRegistered {
		fmt.Println("Broker is already Registered")
		return
	}
	//TODO: Get IP address from Docker ENV
	worker.broker.IP = "66.220.158.68"
	worker.findBrokerDomainName()
	worker.findBrokerGeolocation()

	worker.registerTicker = time.NewTicker(RegisterInterval)
	go func(){
		registrationSuccess := false
		for _ = range worker.registerTicker.C {
			fmt.Println("RegistrationTicker Tick")
			if registrationSuccess {
				worker.registerTicker.Stop()
				break
			}
			registrationSuccess = worker.sendRegistrationRequest()
		}
	}()
}

func (worker *BrokerRegistrationWorker) isBrokerRegistered() bool {
	isBrokerRegistered := true
	broker, err := FindBroker()
	if err != nil {
		isBrokerRegistered = false
	} else if broker.ID == "" {
		isBrokerRegistered = false
	}
	return isBrokerRegistered
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
	location, err := geolocationFetcher.SendGeoLocationRequest(worker.broker.IP)
	if err != nil {
		worker.broker.Geolocation = new(models.Geolocation)
		return
	}
	worker.broker.Geolocation = location
}

func (worker *BrokerRegistrationWorker) sendRegistrationRequest() bool {
	//TODO: get own Broker information
	fmt.Println("Sending Broker Registration Request")

	jsonString, _ := json.Marshal(&worker.broker)

	req, err := http.NewRequest("POST", "http://localhost:8080/rest/brokers", bytes.NewBuffer(jsonString))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		fmt.Println(body)
		return false
	}

	responseBroker := new(models.Broker)

	err = json.Unmarshal(body, responseBroker)
	if err != nil {
		fmt.Println("REGISTER BROKER: Unkown response format")
		return false
	}
	fmt.Println(responseBroker)

	 err = StoreBroker(responseBroker)
	if err != nil {
		return false
	}
	return true
}