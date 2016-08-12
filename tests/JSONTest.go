package main

import (
	"github.com/tkrex/IDS/common/models"
	"encoding/json"
	"fmt"
)

func main() {
	domain := models.NewRealWorldDomain("test")
	json.Marshal(domain)

	xmlByteArray := []byte("<start>")

	var jsonData map[string]*json.RawMessage
	err := json.Unmarshal(xmlByteArray,&jsonData)
	fmt.Print(err)
}