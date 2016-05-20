package providing

import (
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"encoding/json"
	"sync"

	"github.com/tkrex/IDS/domainController/persistence"
)

type DomainInformationRESTProvider struct {
	port string
	providerStarted sync.WaitGroup
	providerStopped sync.WaitGroup

}

func NewIDSGatewayWebInterface(port string) *DomainInformationRESTProvider {
	provider := new(DomainInformationRESTProvider)
	provider.providerStarted.Add(1)
	provider.providerStopped.Add(1)
	go provider.run(port)
	provider.providerStarted.Wait()
	return provider
}

func (provider *DomainInformationRESTProvider) run(port string) {
	fmt.Println("IDSGatewayInterface started")
	 provider.providerStarted.Done()

	router := mux.NewRouter()
	router.HandleFunc("/domainInformation/{domain}", provider.handleDomainInformation).Methods("GET")
	http.ListenAndServe(":" + port, router)
}

func (webInterface *DomainInformationRESTProvider) handleDomainInformation(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	fmt.Println("domain Information Request Received")

	requestParameters := mux.Vars(req)
	domainName := requestParameters["domain"]


	dbDelegate, err := persistence.NewDomainControllerDatabaseWorker()

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	domainInformation, err := dbDelegate.FindDomainInformationByDomainName(domainName)
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



