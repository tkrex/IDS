package gateway

import (
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"encoding/json"
	"sync"
	"github.com/tkrex/IDS/domainController"
	"github.com/tkrex/IDS/common/models"
	"html/template"
)

type IDSGatewayWebInterface struct {
	port string
	providerStarted sync.WaitGroup
	providerStopped sync.WaitGroup

}



func NewIDSGatewayWebInterface(port string) *IDSGatewayWebInterface {
	webInterface := new(IDSGatewayWebInterface)
	webInterface.providerStarted.Add(1)
	webInterface.providerStopped.Add(1)
	go webInterface.run(port)
	webInterface.providerStarted.Wait()
	return webInterface
}

func (webInterface *IDSGatewayWebInterface) run(port string) {
	fmt.Print("IDSGatewayInterface started")
	defer webInterface.providerStarted.Done()

	router := mux.NewRouter()
	router.HandleFunc("/domains/{domainName}", webInterface.handleDomainInformation).Methods("GET")
	router.HandleFunc("/rest/brokers", webInterface.handleBrokers).Methods("GET", "POST")


	http.ListenAndServe(":" + port, router)
}

func (webInterface *IDSGatewayWebInterface) handleDomainInformation(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	fmt.Println("domain Information Request Received")
	domainInformation, err := Fi
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	t, _ := template.ParseFiles("templates/domainInformation.html")
	t.Execute(res, domainInformation)
	fmt.Println(req.RemoteAddr)
	res.Header().Set("Content-Type", "application/json")

	outgoingJSON, error := json.Marshal(domainInformation)

	if error != nil {
		fmt.Println(error.Error())
		http.Error(res, error.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(res, string(outgoingJSON))
}

func (webInterface *IDSGatewayWebInterface) handleBrokers(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	switch req.Method {
	case "GET":
		brokers, err := FindAllBrokers()

		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		outgoingJSON, error := json.Marshal(brokers)

		if error != nil {
			fmt.Println(error.Error())
			http.Error(res, error.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprint(res, string(outgoingJSON))
	case "POST":
		broker := new(models.Broker)
		decoder := json.NewDecoder(req.Body)
		error := decoder.Decode(&broker)
		if error != nil {
			fmt.Println(error.Error())
			http.Error(res, error.Error(), http.StatusInternalServerError)
			return
		}
		registrationHandler := NewBrokerRegistrationHandler()
		response, error  :=registrationHandler.registerBroker(broker)
		if error != nil {
			fmt.Println(error.Error())
			http.Error(res, error.Error(), http.StatusInternalServerError)
			return
		}

		outgoingJSON, error := json.Marshal(response)
		fmt.Println(string(outgoingJSON))
		if error != nil {
			fmt.Println(error.Error())
			http.Error(res, error.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprint(res, string(outgoingJSON))

	}


}


