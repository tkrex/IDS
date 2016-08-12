package main

import (
	"fmt"
	"net/url"
)

func main() {

	string := "http://localhost:8080/domainController/default/new"
	 url,_ := url.Parse(string)
	fmt.Println(url.Path)
}
