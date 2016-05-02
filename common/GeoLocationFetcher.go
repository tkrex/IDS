package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type GeoLocationFetcher struct {
	apiServerAddress string
}


func NewGeoLocationFetcher(apiServerAddress string) *GeoLocationFetcher {
	fetcher := new(GeoLocationFetcher)
	fetcher.apiServerAddress = apiServerAddress
	return fetcher
}

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


func (fetcher *GeoLocationFetcher) sendGeoLocationRequest(address string) (*Location,error)  {
	response, err := http.Get(fetcher.apiServerAddress + address)
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
	return &location,nil
}