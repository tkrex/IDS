package registration


import (
	b64 "encoding/base64"
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
)

//Categorizes a Website using the Web Shrinker API
type WebsiteCategorizationWorker struct {
	accessToken      string
	apiServerAddress string
}

//URL and Authentication Information for Web Shrinker API
const (
	AccessKey = "owLf4fHmY0jMwQLNapZD"
	AccessSecret = "F1ltaStEs68ObV92fTq1"
	Endpoint = "https://api.webshrinker.com/categories/v2/"
)

func NewWebsiteCategorizationWorker() *WebsiteCategorizationWorker {
	worker := new(WebsiteCategorizationWorker)
	return worker
}

//Request and returns categories for a Website URL
func (worker *WebsiteCategorizationWorker) RequestCategoriesForWebsite(websiteAddress string) ([]string, error) {
	base64String := b64.StdEncoding.EncodeToString([]byte(websiteAddress))
	client := http.Client{}

	requestUrl := Endpoint + base64String
	req, _ := http.NewRequest("GET", requestUrl, nil)

	req.SetBasicAuth(AccessKey, AccessSecret)
	req.Header.Add("Host", "api.webshrinker.com")

	response, err := client.Do(req)
	if err != nil {
		fmt.Printf("%s", err)
		return nil, err
	} else {
		defer response.Body.Close()

		if response.StatusCode >= 400 {
			fmt.Println(response.Header)
			return nil, err
		}
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			return nil, err
		}
		fmt.Println(string(contents))
		var jsonData map[string]interface{}
		if err := json.Unmarshal(contents, &jsonData); err != nil {
			fmt.Printf("%s", err)
		}
		data := jsonData["data"].([]interface{})
		innerDict := data[0].(map[string]interface{})
		categoryArray := innerDict["categories"].([]interface{})

		categories := make([]string, len(categoryArray))
		for index, category := range categoryArray {
			categories[index] = category.(string)
		}
		return categories, nil
	}
}