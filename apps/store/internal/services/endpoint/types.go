package endpoint

import (
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
	"github.com/ilya-mezentsev/micro-dep/shared/types/operations"
)

type Service = operations.CUD[models.Endpoint]

type Repo interface {
	operations.CUD[models.Endpoint]
	// Exists - expected that it returns 2 flags: entity existence and endpoint existence
	Exists(model models.Endpoint) (bool, bool, error)
	HasRelation(endpointId models.Id) (bool, error)
}
