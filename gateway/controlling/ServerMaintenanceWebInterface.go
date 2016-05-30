package controlling

import (
	"sync"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/tkrex/IDS/common/models"
	"encoding/json"
	"fmt"
	"net"
)

type ServerMaintenanceWebInterface struct {
	port                             string
	providerStarted                  sync.WaitGroup
	providerStopped                  sync.WaitGroup
}



func NewServerMaintenanceWebInterface(port string) *ServerMaintenanceWebInterface {
	webInterface := new(ServerMaintenanceWebInterface)
	webInterface.providerStarted.Add(1)
	webInterface.providerStopped.Add(1)
	go webInterface.run(port)
	webInterface.providerStarted.Wait()
	fmt.Println("Maintance Web Interface Started")
	return webInterface
}


func (webInterface *ServerMaintenanceWebInterface) run(port string) {
	router := mux.NewRouter()
	router.HandleFunc("/domainController/{domain}/new", webInterface.instantiateDomainController).Methods("GET")
	router.HandleFunc("/domainController/{domain}/delete", webInterface.deleteDomainController).Methods("GET")


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

	requestHandler := NewDomainControllerManagementRequestHandler()
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

	requestHandler := NewDomainControllerManagementRequestHandler()
	if domainController := requestHandler.handleManagementRequest(managementRequest); domainController != nil {
		json.NewEncoder(res).Encode(domainController)
		return
	}
	http.Error(res, "Internal Error", http.StatusInternalServerError)
	return
}