package brokerRegistration

import (
	"sync"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/tkrex/IDS/common/models"
	"encoding/json"
	"fmt"
	"net"
)

type BrokerRegistrationInterface struct {
	port                    string
	providerStarted         sync.WaitGroup
	providerStopped         sync.WaitGroup
}

func NewBrokerRegistrationInterface(port string) *BrokerRegistrationInterface {
	webInterface := new(BrokerRegistrationInterface)
	webInterface.providerStarted.Add(1)
	webInterface.providerStopped.Add(1)
	go webInterface.run(port)
	webInterface.providerStarted.Wait()
	fmt.Println("Maintance Web Interface Started")
	return webInterface
}

func (webInterface *BrokerRegistrationInterface) run(port string) {
	router := mux.NewRouter()
	router.HandleFunc("/rest/brokers", webInterface.addBroker).Methods("POST")

	listener, err := net.Listen("tcp", ":" + port)
	if err != nil {
		panic(err)
	}
	webInterface.providerStarted.Done()
	go http.Serve(listener, router)
}

func (webInterface *BrokerRegistrationInterface) addBroker(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	broker := new(models.Broker)
	decoder := json.NewDecoder(req.Body)
	error := decoder.Decode(&broker)
	if error != nil {
		fmt.Println(error.Error())
		http.Error(res, error.Error(), http.StatusInternalServerError)
		return
	}
	registrationManager := NewBrokerRegistrationManager()
	response, error := registrationManager.RegisterBroker(broker)
	if error != nil {
		fmt.Println(error.Error())
		http.Error(res, error.Error(), http.StatusInternalServerError)
		return
	}

	outgoingJSON, error := json.Marshal(response)
	if error != nil {
		fmt.Println(error.Error())
		http.Error(res, error.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(res, string(outgoingJSON))
}

