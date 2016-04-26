package common
//
//import (
//	"github.com/tkrex/IDS/common/models"
//)
//
//type RequestHandler struct {
//	forwardingTable ForwardingTable
//	pedingRequestTable *PendingRequestTable
//}
//
//type ForwardingTable map[string]string
//type PendingRequestTable chan []*models.Topic
//
//func NewRequestProcessor() *RequestHandler {
//	handler := new(RequestHandler)
//	handler.forwardingTable = new(PendingRequestTable)
//	return handler
//}
//
//func (handler *RequestHandler) HandleTopicsRequest(domain string) (*[]models.Topic,error) {
//	//
//	//pendingRequest, _ := handler.pedingRequestTable[domain]
//	//if pendingRequest != nil {
//	//	topics := pendingRequest
//	//	return (topics, nil)
//	//}
//	//forwardAddress, ok := handler.forwardingTable[domain]
//	//if !ok {
//	//	return errors.New("No matching Domain found")
//	//}
//	//
//	//
//
//}
//
//func (handler *RequestHandler) forwardRequest(domain string, resultChannel chan []*models.Topic) {
//
//}