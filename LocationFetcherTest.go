package main

import (

	"fmt"
	"github.com/tkrex/IDS/daemon/registration"
)

func main() {
	fetcher := registration.NewGeoLocationFetcher()
	location,_ := fetcher.SendGeoLocationRequest("ma-krex-bruegge.in.tum.de")
	fmt.Println(location)


}
