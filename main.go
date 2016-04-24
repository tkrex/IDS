package main

import (
	"github.com/gorilla/mux"

)
import (
	"net/http"
	"github.com/tkrex/IDS/common/models"
	"encoding/json"
	"fmt"
)

func main() {
	//flag.Parse()
	//brokerAddress := "tcp://localhost:1883"
	//desiredTopic  := "#"
	//var subscriber common.InformationProducer
	//subscriber = common.NewMqttSubscriber(brokerAddress,desiredTopic)
	//
	//time.Sleep(time.Second * 60)
	//subscriber.Close()


	//worker := common.NewWebsiteCategorizationWorker("owLf4fHmY0jMwQLNapZD","http://api.webshrinker.com/categories/v2")
	//worker.RequestCategoriesForWebsite("www1.in.tum.de")
	router := mux.NewRouter()
	router.HandleFunc("/topics/{domain}", handleTopics).Methods("GET")
	http.ListenAndServe(":8080", router)
}

func handleTopics(res http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	domain := vars["domain"]
	fmt.Println(domain)
	topics := [5]*models.Topic{}
	for i := 0 ; i< 5; i++ {
		topic := models.NewTopic(i,"test",[]byte{})
		topics[i] = topic
	}

	res.Header().Set("Content-Type", "application/json")

	outgoingJSON, error := json.Marshal(topics)

	if error != nil {
		fmt.Println(error.Error())
		http.Error(res, error.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(res, string(outgoingJSON))
}
