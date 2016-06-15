package models

type DomainInformationRequest struct {
	domain    string
	country   string
	topicName string
}


func NewDomainInformationRequest(domain, country, topicName string) *DomainInformationRequest {
	request := new(DomainInformationRequest)
	request.domain = domain
	request.country = country
	request.topicName = topicName
	return request
}

func (request *DomainInformationRequest) Country() string {
	return request.country
}

func (request *DomainInformationRequest) Domain() string {
	return request.domain
}

func (request *DomainInformationRequest) Name () string {
	return request.topicName
}