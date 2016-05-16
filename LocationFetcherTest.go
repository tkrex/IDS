package main

import (
	"github.com/tkrex/IDS/common/models"
	"encoding/json"
	"net/http"
	"bytes"
	"fmt"
	"io/ioutil"
)

func main() {
	domain1 := models.NewRealWorldDomain("domain1")
 	domainC1 := models.NewDomainController("address1",domain1)
	domains := []*models.DomainController{domainC1}
	json,_:= json.Marshal(domains)


	req, err := http.NewRequest("POST", "http://localhost:8080/controlling", bytes.NewBuffer(json))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

}
