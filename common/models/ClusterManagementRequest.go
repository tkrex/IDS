package models

type ClusterManagementRequestType string
const (
	DomainControllerStart ClusterManagementRequestType = "Update"
	DomainControllerStop ClusterManagementRequestType = "Delete"
	DomainControllerFetch ClusterManagementRequestType = "Fetch"
)

//Contains information of a ClusterManagementRequest
type ClusterManagementRequest struct {
	RequestType  ClusterManagementRequestType
	Domain       *RealWorldDomain
	ParentDomain *RealWorldDomain
}

func NewClusterManagementRequest(requestType ClusterManagementRequestType, domain *RealWorldDomain) *ClusterManagementRequest {
	request := new(ClusterManagementRequest)
	request.RequestType = requestType
	request.Domain = domain
	return request
}
