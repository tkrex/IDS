package models

type DomainControllerManagementRequest struct {
	RequestType  ControlMessageType
	Domain       *RealWorldDomain
	ParentDomain *RealWorldDomain
}

func NewDomainControllerManagementRequest(requestType ControlMessageType, domain *RealWorldDomain) *DomainControllerManagementRequest{
	request := new(DomainControllerManagementRequest)
	request.RequestType = requestType
	request.Domain = domain
	return request
}
