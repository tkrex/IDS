package models

type Geolocation struct {
	Country string `json:"country"`
	Region  string `json:"region"`
	City    string `json:"city"`
}

func NewGeolocation(county, region, city string) *Geolocation {
	location := new(Geolocation)
	location.Country = county
	location.Region = region
	location.City = city
	return location
}