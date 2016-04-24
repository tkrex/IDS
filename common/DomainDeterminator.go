package common

import (
	b64 "encoding/base64"
	"net/http"
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
)

type WebsiteCategorizationWorker struct  {
	accessToken string
	apiServerAddress string
}

func NewWebsiteCategorizationWorker(accessToken string, apiServerAddress string) *WebsiteCategorizationWorker {
	worker := new(WebsiteCategorizationWorker)
	worker.apiServerAddress = apiServerAddress
	worker.accessToken = accessToken
	return  worker
}


func (worker *WebsiteCategorizationWorker) RequestCategoriesForWebsite(websiteAddress string) {
	base64String := b64.StdEncoding.EncodeToString([]byte(websiteAddress))
	client := http.Client{}

	requestUrl := worker.apiServerAddress  +"/"+base64String
	req, _ := http.NewRequest("GET",requestUrl,nil)

	req.SetBasicAuth("owLf4fHmY0jMwQLNapZD","F1ltaStEs68ObV92fTq1")
	req.Header.Add("Host","api.webshrinker.com")
	fmt.Println(req.Header)


	response, err := client.Do(req)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		defer response.Body.Close()

		if response.StatusCode >= 400 {
			fmt.Println(response.Header)
			return
		}
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		var jsonData map[string]interface{}
		if err := json.Unmarshal(contents, &jsonData); err != nil {
			fmt.Printf("%s", err)
		}

		data := jsonData["data"].([]interface{})
		innerDict := data[0].(map[string]interface{})
		categoryArray := innerDict["categories"].([]string)
		fmt.Println(categoryArray)
	}
}