package configuration

import (
	"sync"
	"github.com/tkrex/IDS/common/models"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/tkrex/IDS/daemon/persistence"
	"encoding/json"
	"os"
)

type ConfigurationInterface struct {
	port                           string
	providerStarted                sync.WaitGroup
	providerStopped                sync.WaitGroup
}



func NewConfigurationInterface(port string) *ConfigurationInterface {
	webInterface := new(ConfigurationInterface)
	webInterface.providerStarted.Add(1)
	webInterface.providerStopped.Add(1)
	go webInterface.run(port)
	webInterface.providerStarted.Wait()
	fmt.Println("Daemon Configuration Interface Started")
	return webInterface
}

func (webInterface *ConfigurationInterface) run(port string) {
	webInterface.providerStarted.Done()
	goPath := os.Getenv("GOPATH")
	htmlPath := goPath+"/src/github.com/tkrex/IDS/daemon/frontend/"
	fs := http.Dir(htmlPath)
	fileHandler := http.FileServer(fs)


	router := mux.NewRouter()
	router.Handle("/", http.RedirectHandler("/static/", 302))
	router.HandleFunc("/topics", webInterface.handleTopics).Methods("POST","GET")
	router.HandleFunc("/topics/{name}",webInterface.findTopicByName).Methods("GET")
	router.HandleFunc("/broker",webInterface.findBroker).Methods("GET")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static", fileHandler))
	http.ListenAndServe(":" + port, router)
}



func (webInterface *ConfigurationInterface) findTopicByName(res http.ResponseWriter, req *http.Request) {
	dbDelegate, err := persistence.NewDaemonDatabaseWorker()
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dbDelegate.Close()

	parameters := mux.Vars(req)
	name := parameters["name"]
	topic, err := dbDelegate.FindTopicByName(name)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(res).Encode(topic)
}
func (webInterface ConfigurationInterface) handleTopics(res http.ResponseWriter, req *http.Request) {
	dbDelegate, err := persistence.NewDaemonDatabaseWorker()
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dbDelegate.Close()
	switch req.Method {
	case "POST":
		topic := new(models.Topic)
		decoder := json.NewDecoder(req.Body)
		error := decoder.Decode(&topic)
		if error != nil {
			fmt.Println(error.Error())
			http.Error(res, error.Error(), http.StatusInternalServerError)
			return
		}
		dbDelegate.UpdateTopicDomainAndVisibility(topic)

		fmt.Fprint(res, "OK")
	case "GET":
		topics, err := dbDelegate.FindAllTopics()
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonData, err := json.Marshal(topics)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return

		}
		fmt.Fprint(res, string(jsonData))
	}
}


func (webInterface *ConfigurationInterface) findBroker(res http.ResponseWriter, req *http.Request) {
	dbDelegate, err := persistence.NewDaemonDatabaseWorker()
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dbDelegate.Close()

	broker, err := dbDelegate.FindBroker()
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(res).Encode(broker)
}