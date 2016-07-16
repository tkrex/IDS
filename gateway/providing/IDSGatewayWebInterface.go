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
	"strings"
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
	router.HandleFunc("/rest/domainInformation/{domain}", webInterface.handleDomainInformation).Methods("GET")
	router.HandleFunc("/rest/brokers/{brokerId}/{domain}", webInterface.getDomainInformationForBroker).Methods("GET")
	router.HandleFunc("/rest/domainControllers/{domainName}", webInterface.getDomainControllerForDomain).Methods("GET")
	router.HandleFunc("/rest/brokers/{domainName}", webInterface.getBrokersForDomain).Methods("GET")
	router.HandleFunc("/rest/brokers", webInterface.addBroker).Methods("POST")
	router.HandleFunc("/rest/brokers", webInterface.addBroker).Methods("GET")
	router.HandleFunc("/rest/domains", webInterface.getAllDomains).Methods("GET")

	router.PathPrefix("/").Handler(http.StripPrefix("/", fileHandler))
	http.ListenAndServe(":" + port, router)
}

func (webInterface *IDSGatewayWebInterface) handleDomainInformation(res http.ResponseWriter, req *http.Request) {
	//res.Header().Set("Content-Type", "application/json")
	fmt.Println("domain Information Request Received")
	requestParameters := mux.Vars(req)
	domainName := requestParameters["domain"]
	req.ParseForm()
	location := req.FormValue("location")
	parsedLocation := new(models.Geolocation)
	json.Unmarshal([]byte(location),parsedLocation)

	name := req.FormValue("name")

	informationRequest := models.NewDomainInformationRequest(domainName, parsedLocation,name)

	requestHandler := NewDomainInformationRequestHandler()
	domainInformation := requestHandler.handleRequest(informationRequest)
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
	fmt.Println("Domain Controller Request Received")

	requestParameters := mux.Vars(req)
	domainName := requestParameters["domainName"]
	if domainName == "" {
		http.Error(res, "Error", http.StatusInternalServerError)
		return
	}
	domain := models.NewRealWorldDomain(domainName)
	domainController, _ := routing.NewRoutingManager().DomainControllerForDomain(domain,false)
	if domainController != nil {
		fmt.Println("Responding with Domain Controller: ",domainController)
		json.NewEncoder(res).Encode(domainController)
		return
	}
	http.Error(res,"No DomainController Found",http.StatusNoContent)
}

func (webInterface *IDSGatewayWebInterface) getDomainInformationForBroker(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	fmt.Println("Domain Information Request Received")

	requestParameters := mux.Vars(req)
	brokerId := requestParameters["brokerId"]
	domainName := requestParameters["domain"]
	fmt.Println(domainName)
	req.ParseForm()
	location := req.FormValue("location")
	parsedLocation := new(models.Geolocation)
	json.Unmarshal([]byte(location),parsedLocation)
	name := req.FormValue("name")

	informationRequest := models.NewDomainInformationRequest(domainName, parsedLocation,name)

	domainInformation, err := NewDomainInformationForBrokerRequestHandler().handleRequest(brokerId,informationRequest)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	outgoingJSON, error := json.Marshal(domainInformation)

	if error != nil {
		fmt.Println(error.Error())
		http.Error(res, error.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(res, string(outgoingJSON))
}

func (webInterface *IDSGatewayWebInterface) getBrokersForDomain(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	fmt.Println("Received Broker Request")
	requestParameters := mux.Vars(req)
	domainName := requestParameters["domainName"]
	domainName = strings.Replace(domainName,"%2F","/",-1)
	req.ParseForm()
	location := req.FormValue("location")
	parsedLocation := new(models.Geolocation)
	json.Unmarshal([]byte(location),parsedLocation)
	fmt.Println(parsedLocation)
	name := req.FormValue("name")

	informationRequest := models.NewDomainInformationRequest(domainName,parsedLocation,name)
	brokersSortedByDomains,err := NewBrokerRequestHandler().handleRequest(informationRequest)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	outgoingJSON, error := json.Marshal(brokersSortedByDomains)

	if error != nil {
		fmt.Println(error.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(res, string(outgoingJSON))
}

func (webInterface *IDSGatewayWebInterface) getAllBrokers(res http.ResponseWriter, req *http.Request){
	//TODO: implemenet
}


func (webInterface *IDSGatewayWebInterface) getAllDomains(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Domain Request Received")
	domains , err := NewDomainRequestHandler().handleRequest()
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(res).Encode(domains)
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
	if error != nil {
		fmt.Println(error.Error())
		http.Error(res, error.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(res, string(outgoingJSON))
}



