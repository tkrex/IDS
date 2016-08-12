package registration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"github.com/tkrex/IDS/common/models"
)
//Determines Geolocation of Broker based on IP Address
type GeoLocationFetcher struct {
}


func NewGeoLocationFetcher() *GeoLocationFetcher {
	fetcher := new(GeoLocationFetcher)
	return fetcher
}
//Endpoint of Freegeoip API
const APIEndpoint = "http://freegeoip.net/json/"

//Reponse of API request
type Location struct {
	IP          string  `json:"id"`
	CountryCode string  `json:"country_code"`
	CountryName string  `json:"country_name"`
	RegionCode  string  `json:"region_code"`
	RegionName  string  `json:"region_name"`
	City        string  `json:"city"`
	ZipCode     string  `json:"zip_code"`
	TimeZone    string  `json:"time_zone"`
	Latitude    float32 `json:"latitude"`
	Longitude   float32 `json:"longitude"`
	MetroCode   int     `json:"metro_code"`
}

//Requests and returns geolocation based on IP Address
func (fetcher *GeoLocationFetcher) SendGeoLocationRequest(address string) (*models.Geolocation,error)  {
	response, err := http.Get(APIEndpoint + address)
	if err != nil {
		fmt.Println(err)
		return nil,err

	}
	fmt.Println(response.StatusCode)
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return nil,err
	}
	location := Location{}
	json.Unmarshal([]byte(string(contents)), &location)
	geolocation := models.NewGeolocation(location.CountryName,location.RegionName, location.City,location.Longitude,location.Latitude)
	return geolocation,nil
}
