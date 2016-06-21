package providing

import (
	"github.com/tkrex/IDS/gateway/persistence"
	"errors"
	"github.com/tkrex/IDS/common/models"
)

type DomainRequestHandler struct {

}

func NewDomainRequestHandler() *DomainRequestHandler {
	return new(DomainRequestHandler)
}

func (requestHandler *DomainRequestHandler) handleRequest() ([]*models.RealWorldDomain, error){

	return []*models.RealWorldDomain{ models.NewRealWorldDomain("education/schools"), models.NewRealWorldDomain("education"), models.NewRealWorldDomain("education/university")}, nil
	dbDelegate := persistence.NewGatewayDBWorker()
	if dbDelegate == nil {
		return nil, errors.New("No connection to database")
	}
	domains, error := dbDelegate.FindAllDomains()
	return domains, error
}
