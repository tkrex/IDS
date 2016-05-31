package registration

import (
	"github.com/tkrex/IDS/common/models"
	"encoding/json"
	"net/http"
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"time"
	"sync"
	"github.com/tkrex/IDS/daemon/persistence"
	"github.com/tkrex/IDS/common/controlling"
	"os"
)

const RegisterInterval = time.Second * 10

type BrokerRegistrationWorker struct {
	registrationServerAddress string
	registerTicker *time.Ticker
	workerStarted  sync.WaitGroup
	workerStopped  sync.WaitGroup
	dbDelegate     *persistence.DaemonDatabaseWorker
}

func NewBrokerRegistrationWorker(registrationServerAddress string) *BrokerRegistrationWorker {
	worker := new(BrokerRegistrationWorker)
	worker.registrationServerAddress = registrationServerAddress
	worker.workerStarted.Add(1)
	worker.workerStopped.Add(1)
	go worker.registerBroker()
	worker.workerStarted.Wait()
	return worker
}

func (worker *BrokerRegistrationWorker) registerBroker() {
	databaseWorker, err := persistence.NewDaemonDatabaseWorker()
	if err != nil {
		fmt.Println("Database not reachable")
		return
	}
	worker.dbDelegate = databaseWorker
	worker.workerStarted.Done()

	if isBrokerRegistered := worker.isBrokerRegistered(); isBrokerRegistered {
		fmt.Println("Broker is already Registered")
		return
	}
	//TODO: Get IP address from Docker ENV

	broker := models.NewBroker()
	broker.IP = os.Getenv("BROKER_URI")
	worker.findDomainNameForBroker(broker)
	worker.findRealWorldDomainsForBroker(broker)
	worker.findGeolocationForBroker(broker)

	worker.registerTicker = time.NewTicker(RegisterInterval)
	go func() {
		registrationSuccess := false
		defer worker.dbDelegate.Close()
		for _ = range worker.registerTicker.C {
			fmt.Println("RegistrationTicker Tick")
			if registrationSuccess {
				worker.registerTicker.Stop()
				break
			}
			registrationSuccess = worker.sendRegistrationRequestForBroker(broker)
		}
	}()
}

func (worker *BrokerRegistrationWorker) isBrokerRegistered() bool {
	isBrokerRegistered := true
	broker, err := worker.dbDelegate.FindBroker()
	if err != nil {
		isBrokerRegistered = false
	} else if broker.ID == "" {
		isBrokerRegistered = false
	}
	return isBrokerRegistered
}

func (worker *BrokerRegistrationWorker) findDomainNameForBroker(broker *models.Broker) {
	name, err := net.LookupAddr(broker.IP)
	if err != nil {
		broker.InternetDomain = ""
		return
	}
	broker.InternetDomain = name[0]
}

func (worker *BrokerRegistrationWorker) findRealWorldDomainsForBroker(broker *models.Broker) {
	categorizer := NewWebsiteCategorizationWorker()
	categories,_ := categorizer.RequestCategoriesForWebsite("www.in.tum.de")
	broker.RealWorldDomains = make([]*models.RealWorldDomain,len(categories))
	for index,category := range categories {
		domain := models.NewRealWorldDomain(category)
		broker.RealWorldDomains[index] = domain
	}
}

func (worker *BrokerRegistrationWorker) findGeolocationForBroker(broker *models.Broker) {
	geolocationFetcher := NewGeoLocationFetcher()
	location, err := geolocationFetcher.SendGeoLocationRequest(broker.IP)
	if err != nil {
		broker.Geolocation = new(models.Geolocation)
		return
	}
	broker.Geolocation = location
}

func (worker *BrokerRegistrationWorker) sendRegistrationRequestForBroker(broker *models.Broker) bool {
	fmt.Println("Sending Broker Registration Request")

	jsonString, _ := json.Marshal(&broker)

	req, err := http.NewRequest("POST", worker.registrationServerAddress+"/rest/brokers", bytes.NewBuffer(jsonString))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		fmt.Println(resp.Status)
		fmt.Println(body)
		return false
	}

	response := new(models.BrokerRegistrationResponse)

	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("REGISTER BROKER: Unkown response format")
		return false
	}

	err = worker.dbDelegate.StoreBroker(response.Broker)

	controllingDBDelegate, error := controlling.NewControlMessageDBDelegate()
	if error != nil {
		fmt.Println(error)
		return false
	}
	defer controllingDBDelegate.Close()
	err = controllingDBDelegate.StoreDomainControllers(response.DomainControllers)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}