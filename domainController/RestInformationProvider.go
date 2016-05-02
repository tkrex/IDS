package domainController

import (
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"encoding/json"
	"sync"
)

type RestInformationProvider struct {
	persistenceManager DomainInformationPersistenceManager
	port string
	providerStarted sync.WaitGroup
	providerStopped sync.WaitGroup

}


func NewRestInformationProvider(persistenceManager DomainInformationPersistenceManager, port string) *RestInformationProvider {
	provider := new(RestInformationProvider)
	provider.persistenceManager = persistenceManager
	provider.providerStarted.Add(1)
	provider.providerStopped.Add(1)
	go provider.run(port)
	provider.providerStarted.Wait()
	return provider
}

func (provider *RestInformationProvider) run(port string) {
	defer func(){
		provider.providerStarted.Done()
	}()
	router := mux.NewRouter()
	router.HandleFunc("/domains/{domainName}", provider.handleDomainInformation).Methods("GET")
	router.HandleFunc("/brokers", provider.handleBrokers).Methods("GET")

	http.ListenAndServe(":" + port, router)
}

func (provider *RestInformationProvider) handleDomainInformation(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	requestProcessor := NewRequestProcessor()
	domainInformation ,error := requestProcessor.handleDomainInformationRequest(req)
	if error != nil {
		fmt.Println(error.Error())
		http.Error(res, error.Error(), http.StatusInternalServerError)
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

func (provider *RestInformationProvider) handleBrokers(res http.ResponseWriter, req *http.Request) {
	brokers := provider.persistenceManager.Brokers()
	fmt.Println(req.RemoteAddr)
	res.Header().Set("Content-Type", "application/json")

	outgoingJSON, error := json.Marshal(brokers)

	if error != nil {
		fmt.Println(error.Error())
		http.Error(res, error.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(res, string(outgoingJSON))
}

