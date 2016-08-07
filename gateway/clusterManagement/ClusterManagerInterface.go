package clusterManagement

import (
	"sync"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/tkrex/IDS/common/models"
	"encoding/json"
	"fmt"
	"net"
)

type ClusterManagerInterface struct {
	port                    string
	providerStarted         sync.WaitGroup
	providerStopped         sync.WaitGroup
}

func NewClusterManagerInterface(port string) *ClusterManagerInterface {
	webInterface := new(ClusterManagerInterface)
	webInterface.providerStarted.Add(1)
	webInterface.providerStopped.Add(1)
	go webInterface.run(port)
	webInterface.providerStarted.Wait()
	fmt.Println("Maintance Web Interface Started")
	return webInterface
}

func (webInterface *ClusterManagerInterface) run(port string) {
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

func (webInterface ClusterManagerInterface) instantiateDomainController(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "*")

	parameters := mux.Vars(req)
	domainName := parameters["domain"]
	domain := models.NewRealWorldDomain(domainName)
	req.ParseForm()
	parentDomain := req.FormValue("parent_domain")
	fmt.Println(parentDomain)

	var messageType = models.DomainControllerStart
	managementRequest := models.NewClusterManagementRequest(messageType, domain)
	managementRequest.ParentDomain = models.NewRealWorldDomain(parentDomain)

	manager := NewClusterManager()
	domainController, error := manager.HandleManagementRequest(managementRequest)
	if error != nil {
		http.Error(res, error.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(res).Encode(&domainController)
}

func (webInterface ClusterManagerInterface) fetchDomainController(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "*")

	parameters := mux.Vars(req)
	domainName := parameters["domain"]
	domain := models.NewRealWorldDomain(domainName)

	var messageType = models.DomainControllerFetch
	managementRequest := models.NewClusterManagementRequest(messageType, domain)

	requestHandler := NewClusterManager()
	domainController, error := requestHandler.HandleManagementRequest(managementRequest)
	if error != nil {
		http.Error(res, error.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(res).Encode(&domainController)
}

func (webInterface ClusterManagerInterface) deleteDomainController(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	parameters := mux.Vars(req)
	domainName := parameters["domain"]
	domain := models.NewRealWorldDomain(domainName)

	var messageType = models.DomainControllerStop
	managementRequest := models.NewClusterManagementRequest(messageType, domain)

	requestHandler := NewClusterManager()
	if _, error := requestHandler.HandleManagementRequest(managementRequest); error != nil {
		http.Error(res, error.Error(), http.StatusInternalServerError)
		return
	}
}