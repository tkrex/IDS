package models

type Geolocation struct {
	Country string `json:"country"`
	Region  string `json:"region"`
	City    string `json:"city"`
	Longitude float32 `json:"longitude"`
	Latitude float32 `json:"latitude"`
}

func NewGeolocation(country, region, city string, long float32,lat float32) *Geolocation {
	location := new(Geolocation)
	location.Country = country
	location.Region = region
	location.City = city
	location.Longitude = long
	location.Latitude = lat
	return location
}