package configuration

import (
	"sync"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/tkrex/IDS/common/models"
	"encoding/json"
)

type DomainControllerAllocationInterface struct {
	port                           string
	providerStarted                sync.WaitGroup
	providerStopped                sync.WaitGroup
}



func NewClusterControllerInterface(port string) *DomainControllerAllocationInterface {
	webInterface := new(DomainControllerAllocationInterface)
	webInterface.providerStarted.Add(1)
	webInterface.providerStopped.Add(1)
	go webInterface.run(port)
	webInterface.providerStarted.Wait()
	fmt.Println("Cluster Controller Interface Started")
	return webInterface
}

func (webInterface *DomainControllerAllocationInterface) run(port string) {
	webInterface.providerStarted.Done()


	router := mux.NewRouter()
	router.HandleFunc("/domainControllerAllocation/{domain}", webInterface.handleDomainControllerAllocationRequest).Methods("GET")
	http.ListenAndServe(":" + port, router)
}

func (webInterface *DomainControllerAllocationInterface) handleDomainControllerAllocationRequest(res http.ResponseWriter, req *http.Request) {
	parameters := mux.Vars(req)
	domainName := parameters["domain"]
	domainController , err := webInterface.startNewDomainControllerInstance(domainName)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(res).Encode(domainController)
}


func (webInterface *DomainControllerAllocationInterface) startNewDomainControllerInstance(domainName string) (*models.DomainController, error) {
	domainController := models.NewDomainController("localhost:1883",models.NewRealWorldDomain(domainName))
	return domainController
}



