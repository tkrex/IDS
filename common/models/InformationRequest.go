package models

type DomainInformationRequest struct {
	domain    string
	location  *Geolocation
	topicName string
}


func NewDomainInformationRequest(domain string, location *Geolocation, topicName string) *DomainInformationRequest {
	request := new(DomainInformationRequest)
	request.domain = domain
	request.location = location
	request.topicName = topicName
	return request
}

func (request *DomainInformationRequest) Location() *Geolocation {
	return request.location
}

func (request *DomainInformationRequest) Domain() string {
	return request.domain
}

func (request *DomainInformationRequest) Name () string {
	return request.topicName
}