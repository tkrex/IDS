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

type ServerMaintenanceWebInterface struct {
	port                             string
	managementBrokerAddress		*url.URL
	providerStarted                  sync.WaitGroup
	providerStopped                  sync.WaitGroup
}



func NewServerMaintenanceWebInterface(port string, managementBrokerAddress *url.URL) *ServerMaintenanceWebInterface {
	webInterface := new(ServerMaintenanceWebInterface)
	webInterface.managementBrokerAddress = managementBrokerAddress
	webInterface.providerStarted.Add(1)
	webInterface.providerStopped.Add(1)
	go webInterface.run(port)
	webInterface.providerStarted.Wait()
	fmt.Println("Maintance Web Interface Started")
	return webInterface
}


func (webInterface *ServerMaintenanceWebInterface) run(port string) {
	router := mux.NewRouter()
	router.HandleFunc("/rest/domainControllers/{domain}/new", webInterface.instantiateDomainController).Methods("GET")
	router.HandleFunc("/rest/domainControllers/{domain}/delete", webInterface.deleteDomainController).Methods("GET")
	router.HandleFunc("/rest/domainControllers/{domain}", webInterface.fetchDomainController).Methods("GET")


	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}
	webInterface.providerStarted.Done()
	go http.Serve(listener, router)
}

func (webInterface ServerMaintenanceWebInterface) instantiateDomainController(res http.ResponseWriter, req *http.Request) {
	parameters := mux.Vars(req)
	domainName := parameters["domain"]
	domain := models.NewRealWorldDomain(domainName)

	var messageType = models.DomainControllerChange
	managementRequest := models.NewDomainControllerManagementRequest(messageType,domain)

	requestHandler := NewDomainControllerManagementRequestHandler(webInterface.managementBrokerAddress)
	if domainController := requestHandler.handleManagementRequest(managementRequest); domainController != nil {
		json.NewEncoder(res).Encode(&domainController)
		return
	}
	http.Error(res, "Internal Error", http.StatusInternalServerError)
	return
}

func (webInterface ServerMaintenanceWebInterface) fetchDomainController(res http.ResponseWriter, req *http.Request) {
	parameters := mux.Vars(req)
	domainName := parameters["domain"]
	domain := models.NewRealWorldDomain(domainName)

	var messageType = models.DomainControllerFetch
	managementRequest := models.NewDomainControllerManagementRequest(messageType,domain)

	requestHandler := NewDomainControllerManagementRequestHandler(webInterface.managementBrokerAddress)
	if domainController := requestHandler.handleManagementRequest(managementRequest); domainController != nil {
		json.NewEncoder(res).Encode(&domainController)
		return
	}
	http.Error(res, "Internal Error", http.StatusInternalServerError)
	return
}

func (webInterface ServerMaintenanceWebInterface) deleteDomainController(res http.ResponseWriter, req *http.Request) {
	parameters := mux.Vars(req)
	domainName := parameters["domain"]
	domain := models.NewRealWorldDomain(domainName)

	var messageType = models.DomainControllerDelete
	managementRequest := models.NewDomainControllerManagementRequest(messageType,domain)

	requestHandler := NewDomainControllerManagementRequestHandler(webInterface.managementBrokerAddress)
	if domainController := requestHandler.handleManagementRequest(managementRequest); domainController != nil {
		json.NewEncoder(res).Encode(domainController)
		return
	}
	http.Error(res, "Internal Error", http.StatusInternalServerError)
	return
}