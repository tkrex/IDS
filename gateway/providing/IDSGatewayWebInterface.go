package providing

import (
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"encoding/json"
	"sync"
	"github.com/tkrex/IDS/common/models"
	"github.com/tkrex/IDS/gateway/registration"
	"github.com/tkrex/IDS/common/routing"
)

type IDSGatewayWebInterface struct {
	port            string
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
	fs := http.Dir("./gateway/frontend/")
	fileHandler := http.FileServer(fs)
	router.HandleFunc("/rest/domainInformation/{domainName}", webInterface.handleDomainInformation).Methods("GET")
	router.HandleFunc("/rest/brokers/{brokerId}/domainInformation", webInterface.getDomainInformationForBroker).Methods("GET")
	router.HandleFunc("/rest/domainControllers/{domainName}", webInterface.getDomainControllerForDomain).Methods("GET")
	router.HandleFunc("/rest/brokers/{domainName}", webInterface.getBrokers).Methods("GET")
	router.HandleFunc("/rest/brokers", webInterface.addBroker).Methods("POST")
	router.PathPrefix("/").Handler(http.StripPrefix("/", fileHandler))

	http.ListenAndServe(":" + port, router)
}

func (webInterface *IDSGatewayWebInterface) handleDomainInformation(res http.ResponseWriter, req *http.Request) {
	//res.Header().Set("Content-Type", "application/json")
	fmt.Println("domain Information Request Received")
	requestParameters := mux.Vars(req)
	domainName := requestParameters["domainName"]

	requestHandler := NewDomainInformationRequestHandler()
	domainInformation := requestHandler.handleRequest(domainName)
	if domainInformation == nil {
		http.Error(res, "Error", http.StatusInternalServerError)
		return
	}
	outgoingJSON, error := json.Marshal(domainInformation)

	if error != nil {
		fmt.Println(error.Error())
		http.Error(res, "Error", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(res, string(outgoingJSON))
}

func (webInterface *IDSGatewayWebInterface) getDomainControllerForDomain(res http.ResponseWriter, req *http.Request) {
	fmt.Println("domain controller Request Received")

	requestParameters := mux.Vars(req)
	domainName := requestParameters["domainName"]
	if domainName == "" {
		http.Error(res, "Error", http.StatusInternalServerError)
		return
	}
	domain := models.NewRealWorldDomain(domainName)
	domainController := routing.NewRoutingManager().DomainControllerForDomain(domain)
	if domainController != nil {
		fmt.Println("Responding with Domain Controller: ",domainController)
		json.NewEncoder(res).Encode(domainController)
		return
	}
	http.Error(res,"No DomainController Found",http.StatusNoContent)
}

func (webInterface *IDSGatewayWebInterface) getDomainInformationForBroker(res http.ResponseWriter, req *http.Request) {
	//res.Header().Set("Content-Type", "application/json")
	fmt.Println("domain Information Request Received")

	requestParameters := mux.Vars(req)
	brokerId := requestParameters["brokerId"]
	requestHandler := NewDomainInformationForBrokerRequestHandler()
	domainInformation := requestHandler.handleRequest(brokerId)
	if domainInformation == nil {
		http.Error(res, "Error", http.StatusInternalServerError)
		return
	}
	outgoingJSON, error := json.Marshal(domainInformation)

	if error != nil {
		fmt.Println(error.Error())
		http.Error(res, "Error", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(res, string(outgoingJSON))
}

func (webInterface *IDSGatewayWebInterface) getBrokers(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	fmt.Println("Received Broker Request")
	requestParameters := mux.Vars(req)
	domainName := requestParameters["domainName"]

	req.ParseForm()
	fmt.Println(req.Form)
	country := req.FormValue("country")
	region := req.FormValue("region")
	city := req.FormValue("city")
	location := models.NewGeolocation(country,region,city,0,0)
	fmt.Println(location)
	requestHandler := NewBrokerRequestHandler()
	brokers := requestHandler.handleRequest(domainName)

	if brokers == nil {
		http.Error(res, "Error", http.StatusInternalServerError)
		return
	}
	outgoingJSON, error := json.Marshal(brokers)

	if error != nil {
		fmt.Println(error.Error())
		http.Error(res, "Error", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(res, string(outgoingJSON))
}

func (webInterface *IDSGatewayWebInterface) addBroker(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	broker := new(models.Broker)
	decoder := json.NewDecoder(req.Body)
	error := decoder.Decode(&broker)
	if error != nil {
		fmt.Println(error.Error())
		http.Error(res, error.Error(), http.StatusInternalServerError)
		return
	}
	registrationHandler := registration.NewBrokerRegistrationHandler()
	response, error := registrationHandler.RegisterBroker(broker)
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



