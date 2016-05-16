package domainController

import (
	"github.com/tkrex/IDS/common/models"
	"errors"
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type RequestHandler struct {
	forwardingTable map[string]string
}

type PendingRequestTable chan []*models.Topic

func NewRequestProcessor() *RequestHandler {
	handler := new(RequestHandler)
	handler.forwardingTable = ForwardingTable()
	handler.forwardingTable["testDomain"] = "self"

	return handler
}

func (handler *RequestHandler) handleDomainInformationRequest(request *http.Request) ([]*models.DomainInformationMessage,error) {
	vars := mux.Vars(request)
	domain, _ := vars["domainName"]

	forwardAddress, ok := handler.forwardingTable[domain]
	if !ok {
		return nil, errors.New("No matching Domain found")
	}

	if forwardAddress == "self" {
		//Request Information from persistence Layer
		domainInformation, error := FindDomainInformationByDomainName(domain)

		return domainInformation, error

	} else {
		//Forward address to specified address
		response , error := handler.forwardRequest(request,forwardAddress)
		if error != nil {
			return nil, error
		}
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			return nil,err
		}
		var domainInformationMessage []*models.DomainInformationMessage
		json.Unmarshal([]byte(string(contents)), &domainInformationMessage)
		return domainInformationMessage,nil
	}
}

func (handler *RequestHandler) forwardRequest(request *http.Request, newAddress string) (*http.Response, error) {
	client := http.Client{}
	request.RemoteAddr = newAddress
	return client.Do(request)
}