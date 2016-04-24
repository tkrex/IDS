package common

import (
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"encoding/json"
)

type RestInformationProvider struct {
	persistenceManager  InformationPersistenceManager
}

func (provider *RestInformationProvider) StartListing(port string) {

	provider.persistenceManager = NewMemoryPersistenceManager()
	router := mux.NewRouter()
	router.HandleFunc("/topics/{domain}", provider.handleTopics).Methods("GET")
	http.ListenAndServe(":" + port, router)
}

func (provider *RestInformationProvider) handleTopics(res http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	domain := vars["domain"]
	fmt.Println(domain)
	topics := provider.persistenceManager.Topics()
	res.Header().Set("Content-Type", "application/json")

	outgoingJSON, error := json.Marshal(topics)

	if error != nil {
		fmt.Println(error.Error())
		http.Error(res, error.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(res, string(outgoingJSON))
}

