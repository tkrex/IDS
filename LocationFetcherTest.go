package main

import (

	"fmt"
	"github.com/tkrex/IDS/daemon/layers"
)

func main() {
	//fetcher := common.NewGeoLocationFetcher()
	//location,_ := fetcher.SendGeoLocationRequest("ma-krex-bruegge.in.tum.de")
	//fmt.Println(location)

	categorizer := layers.NewWebsiteCategorizationWorker()
	categories,_ :=categorizer.RequestCategoriesForWebsite("www.in.tum.de")
	fmt.Println(categories)

}
