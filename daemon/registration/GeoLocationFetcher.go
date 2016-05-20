package registration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"github.com/tkrex/IDS/common/models"
)

type GeoLocationFetcher struct {
}


func NewGeoLocationFetcher() *GeoLocationFetcher {
	fetcher := new(GeoLocationFetcher)
	return fetcher
}

const APIEndpont = "http://192.168.99.100:8080/json/"
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


func (fetcher *GeoLocationFetcher) SendGeoLocationRequest(address string) (*models.Geolocation,error)  {
	response, err := http.Get(APIEndpont + address)
	if err != nil {
		fmt.Printf("%s", err)
		return nil,err

	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("%s", err)
		return nil,err
	}
	location := Location{}
	json.Unmarshal([]byte(string(contents)), &location)
	geolocation := models.NewGeolocation(location.CountryName,location.RegionName, location.City)
	return geolocation,nil
}
