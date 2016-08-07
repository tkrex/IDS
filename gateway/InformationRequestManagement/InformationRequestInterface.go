package informationRequestManagement

import (
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"encoding/json"
	"sync"
	"github.com/tkrex/IDS/common/models"
	"strings"
	"os"
)

type InformationRequestInterface struct {
	port            string
	providerStarted sync.WaitGroup
	providerStopped sync.WaitGroup
}

func NewInformationRequestInterface(port string) *InformationRequestInterface {
	webInterface := new(InformationRequestInterface)
	webInterface.providerStarted.Add(1)
	webInterface.providerStopped.Add(1)
	go webInterface.run(port)
	webInterface.providerStarted.Wait()
	return webInterface
}

func (webInterface *InformationRequestInterface) run(port string) {
	fmt.Println("IDSGatewayInterface started")
	webInterface.providerStarted.Done()

	router := mux.NewRouter()
	goPath := os.Getenv("GOPATH")
	htmlPath := goPath+"/src/github.com/tkrex/IDS/gateway/frontend/"
	fs := http.Dir(htmlPath)
	fileHandler := http.FileServer(fs)
	router.HandleFunc("/rest/domainInformation/{domain}", webInterface.handleDomainInformation).Methods("GET")
	router.HandleFunc("/rest/brokers/{brokerId}/{domain}", webInterface.getDomainInformationForBroker).Methods("GET")
	router.HandleFunc("/rest/brokers/{domainName}", webInterface.getBrokersForDomain).Methods("GET")
	router.HandleFunc("/rest/domains", webInterface.getAllDomains).Methods("GET")

	router.PathPrefix("/").Handler(http.StripPrefix("/", fileHandler))
	http.ListenAndServe(":" + port, router)
}

func (webInterface *InformationRequestInterface) handleDomainInformation(res http.ResponseWriter, req *http.Request) {
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
	domainInformation := requestHandler.HandleRequest(informationRequest)
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

func (webInterface *InformationRequestInterface) getDomainInformationForBroker(res http.ResponseWriter, req *http.Request) {
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

	domainInformation, err := NewDomainInformationForBrokerRequestHandler().HandleRequest(brokerId,informationRequest)
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

func (webInterface *InformationRequestInterface) getBrokersForDomain(res http.ResponseWriter, req *http.Request) {
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
	brokersSortedByDomains,err := NewBrokerRequestHandler().HandleRequest(informationRequest)

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


func (webInterface *InformationRequestInterface) getAllDomains(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Domain Request Received")
	domains , err := brokerRegistration.NewBrokerRegistrationManager().AvailableDomains()
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(res).Encode(domains)
}






