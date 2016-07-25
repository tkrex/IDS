package providing

import (
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"encoding/json"
	"sync"
	"github.com/tkrex/IDS/common/models"
)

type DomainInformationRESTProvider struct {
	port string
	providerStarted sync.WaitGroup
	providerStopped sync.WaitGroup

}

func NewDomainInformationRESTProvider(port string) *DomainInformationRESTProvider {
	provider := new(DomainInformationRESTProvider)
	provider.providerStarted.Add(1)
	provider.providerStopped.Add(1)
	go provider.run(port)
	provider.providerStarted.Wait()
	return provider
}

func (provider *DomainInformationRESTProvider) run(port string) {
	fmt.Println("DomainController Rest Interface started")
	 provider.providerStarted.Done()

	router := mux.NewRouter()
	router.HandleFunc("/rest/domainInformation/{domain}", provider.handleDomainInformation).Methods("GET")
	router.HandleFunc("/rest/brokers/{domain}", provider.getBrokersForDomain).Methods("GET")
	router.HandleFunc("/rest/brokers/{brokerId}/{domain}", provider.getDomainInformationForBroker).Methods("GET")

	http.ListenAndServe(":" + port, router)
}

func (webInterface *DomainInformationRESTProvider) getBrokersForDomain(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	fmt.Println("Broker Request Received")

	requestParameters := mux.Vars(req)
	domainName := requestParameters["domain"]
	req.ParseForm()
	location := req.FormValue("location")
	name := req.FormValue("name")
	parsedLocation := new(models.Geolocation)
	json.Unmarshal([]byte(location),parsedLocation)

	informationRequest := models.NewDomainInformationRequest(domainName, parsedLocation,name)
	brokers, err := NewBrokerRequestHandler().handleRequest(informationRequest)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
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


func (webInterface *DomainInformationRESTProvider) getDomainInformationForBroker(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	fmt.Println("domain Information Request Received")

	requestParameters := mux.Vars(req)
	domainName := requestParameters["domain"]
	brokerId := requestParameters["brokerId"]
	req.ParseForm()
	location := req.FormValue("location")
	name := req.FormValue("name")
	parsedLocation := new(models.Geolocation)
	json.Unmarshal([]byte(location),parsedLocation)
	informationRequest := models.NewDomainInformationRequest(domainName,parsedLocation,name)

	domainInformation, err := NewDomainInformationForBrokerRequestHandler().handleRequest(informationRequest,brokerId)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(domainInformation)
	outgoingJSON, error := json.Marshal(domainInformation)

	if err != nil {
		fmt.Println(error.Error())
		http.Error(res, "Error", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(res, string(outgoingJSON))
}

func (webInterface *DomainInformationRESTProvider) handleDomainInformation(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	fmt.Println("domain Information Request Received")

	requestParameters := mux.Vars(req)
	domainName := requestParameters["domain"]
	location := req.FormValue("location")
	parsedLocation := new(models.Geolocation)
	json.Unmarshal([]byte(location),parsedLocation)
	name := req.FormValue("name")
	informationRequest := models.NewDomainInformationRequest(domainName,parsedLocation,name)

	domainInformation, err := NewDomainInformationRequestHandler().handleRequest(informationRequest)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	outgoingJSON, error := json.Marshal(domainInformation)


	if err != nil {
		fmt.Println(error.Error())
		http.Error(res, "Error", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(res, string(outgoingJSON))
}



