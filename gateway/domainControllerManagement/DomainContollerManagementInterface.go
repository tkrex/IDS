package controlling

import (
	"sync"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/tkrex/IDS/common/models"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
)

type DomainContollerManagementInterface struct {
	port                    string
	managementBrokerAddress *url.URL
	providerStarted         sync.WaitGroup
	providerStopped         sync.WaitGroup
}

func NewDomainContollerManagementInterface(port string, managementBrokerAddress *url.URL) *DomainContollerManagementInterface {
	webInterface := new(DomainContollerManagementInterface)
	webInterface.managementBrokerAddress = managementBrokerAddress
	webInterface.providerStarted.Add(1)
	webInterface.providerStopped.Add(1)
	go webInterface.run(port)
	webInterface.providerStarted.Wait()
	fmt.Println("Maintance Web Interface Started")
	return webInterface
}

func (webInterface *DomainContollerManagementInterface) run(port string) {
	router := mux.NewRouter()
	router.HandleFunc("/rest/domainControllers/{domain}/new", webInterface.instantiateDomainController).Methods("GET")
	router.HandleFunc("/rest/domainControllers/{domain}/delete", webInterface.deleteDomainController).Methods("GET")
	router.HandleFunc("/rest/domainControllers/{domain}", webInterface.fetchDomainController).Methods("GET")

	listener, err := net.Listen("tcp", ":" + port)
	if err != nil {
		panic(err)
	}
	webInterface.providerStarted.Done()
	go http.Serve(listener, router)
}

func (webInterface DomainContollerManagementInterface) instantiateDomainController(res http.ResponseWriter, req *http.Request) {
	parameters := mux.Vars(req)
	domainName := parameters["domain"]
	domain := models.NewRealWorldDomain(domainName)
	req.ParseForm()
	parentDomain := req.FormValue("parent_domain")
	fmt.Println(parentDomain)

	var messageType = models.DomainControllerStart
	managementRequest := models.NewDomainControllerManagementRequest(messageType, domain)
	managementRequest.ParentDomain = models.NewRealWorldDomain(parentDomain)

	manager := NewDomainControllerManager()
	domainController, error := manager.handleManagementRequest(managementRequest)
	if error != nil {
		http.Error(res, error.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(res).Encode(&domainController)
}

func (webInterface DomainContollerManagementInterface) fetchDomainController(res http.ResponseWriter, req *http.Request) {
	parameters := mux.Vars(req)
	domainName := parameters["domain"]
	domain := models.NewRealWorldDomain(domainName)

	var messageType = models.DomainControllerFetch
	managementRequest := models.NewDomainControllerManagementRequest(messageType, domain)

	requestHandler := NewDomainControllerManager()
	domainController, error := requestHandler.handleManagementRequest(managementRequest)
	if error != nil {
		http.Error(res, error.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(res).Encode(&domainController)
}

func (webInterface DomainContollerManagementInterface) deleteDomainController(res http.ResponseWriter, req *http.Request) {
	parameters := mux.Vars(req)
	domainName := parameters["domain"]
	domain := models.NewRealWorldDomain(domainName)

	var messageType = models.DomainControllerStop
	managementRequest := models.NewDomainControllerManagementRequest(messageType, domain)

	requestHandler := NewDomainControllerManager()
	if _, error := requestHandler.handleManagementRequest(managementRequest); error != nil {
		http.Error(res, error.Error(), http.StatusInternalServerError)
		return
	}
}