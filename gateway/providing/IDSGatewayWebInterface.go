package providing

import (
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"encoding/json"
	"sync"
	"github.com/tkrex/IDS/common/models"
	"github.com/tkrex/IDS/gateway/persistence"
	"github.com/tkrex/IDS/gateway/registration"
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
	fmt.Println("IDSGatewayInterface started")
	 webInterface.providerStarted.Done()

	router := mux.NewRouter()
	router.HandleFunc("/domains/{domainName}", webInterface.handleDomainInformation).Methods("GET")
	router.HandleFunc("/rest/brokers", webInterface.handleBrokers).Methods("GET", "POST")


	http.ListenAndServe(":" + port, router)
}

func (webInterface *IDSGatewayWebInterface) handleDomainInformation(res http.ResponseWriter, req *http.Request) {
	//res.Header().Set("Content-Type", "application/json")
	fmt.Println("domain Information Request Received")

	requestParameters := mux.Vars(req)
	domainName := requestParameters["domain"]
	requestHandler := NewDomainInformationRequestHandler()
	domainInformation := requestHandler.handleRequest(domainName)
	if domainInformation == nil {
		http.Error(res, "Error", http.StatusInternalServerError)
		return
	}

	t := template.New("domainInformation.html")
	var parseError error


	t, parseError = t.ParseFiles("gateway/templates/domainInformation.html")
	if parseError != nil {
		http.Error(res, parseError.Error(), http.StatusInternalServerError)
		return
	}
	executeError := t.Execute(res,domainInformation)
	if executeError != nil {
		http.Error(res, executeError.Error(), http.StatusInternalServerError)
		return
	}



	//outgoingJSON, error := json.Marshal(domainInformation)
	//
	//if error != nil {
	//	fmt.Println(error.Error())
	//	http.Error(res, "Error", http.StatusInternalServerError)
	//	return
	//}
	//fmt.Fprint(res, string(outgoingJSON))
}

func (webInterface *IDSGatewayWebInterface) handleBrokers(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	switch req.Method {
	case "GET":
		fmt.Println("Received Broker Request")
		dbWorker := persistence.NewGatewayDBWorker()
		if dbWorker == nil {
			fmt.Println("Can't connect to database")
			return
		}
		brokers, err := dbWorker.FindAllBrokers()

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
		registrationHandler := registration.NewBrokerRegistrationHandler()
		response, error  := registrationHandler.RegisterBroker(broker)
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


