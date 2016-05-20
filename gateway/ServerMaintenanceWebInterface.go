package gateway

import (
	"sync"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/tkrex/IDS/common/models"
	"encoding/json"
	"fmt"
)

type ServerMaintenanceWebInterface struct {
	port                           string
	providerStarted                sync.WaitGroup
	providerStopped                sync.WaitGroup
	incomingControlMessagesChannel chan *models.ControlMessage
}



func NewServerMaintenanceWebInterface(port string) *ServerMaintenanceWebInterface {
	webInterface := new(ServerMaintenanceWebInterface)
	webInterface.incomingControlMessagesChannel = make(chan *models.ControlMessage,1000)
	webInterface.providerStarted.Add(1)
	webInterface.providerStopped.Add(1)
	go webInterface.run(port)
	webInterface.providerStarted.Wait()
	fmt.Println("Maintance Web Interface Started")
	return webInterface
}

func (webInterface *ServerMaintenanceWebInterface) IncomingControlMessagesChannel() chan []*models.DomainController {
	return webInterface.incomingControlMessagesChannel
}

func (webInterface *ServerMaintenanceWebInterface) run(port string) {
	webInterface.providerStarted.Done()
	router := mux.NewRouter()
	router.HandleFunc("/controlling", webInterface.handleControlMessages).Methods("POST","DELETE")
	http.ListenAndServe(":" + port, router)
}



func (webInterface ServerMaintenanceWebInterface) handleControlMessages(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Received control POST request")
	domainControllerInformation := []*models.DomainController{}
	decoder := json.NewDecoder(req.Body)
	error := decoder.Decode(&domainControllerInformation)
	if error != nil {
		fmt.Println(error.Error())
		http.Error(res, error.Error(), http.StatusInternalServerError)
		return
	}

	controlMessage := new(models.ControlMessage)
	switch req.Method {
	case "POST":
		controlMessage.MessageType = models.DomainControllerUpdate
	case "DELETE":
		controlMessage.MessageType = models.DomainControllerDelete
	}


	fmt.Fprint(res, nil)
	webInterface.incomingControlMessagesChannel <- controlMessage

}