package endpoint

import (
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
	"github.com/ilya-mezentsev/micro-dep/shared/types/operations"
	"github.com/ilya-mezentsev/micro-dep/store/internal/services/shared"
)

type Service = operations.CUD[shared.Endpoint]

type Repo interface {
	operations.CUD[shared.Endpoint]
	// Exists - expected that it returns 2 flags: entity existence and endpoint existence
	Exists(model shared.Endpoint) (bool, bool, error)
	HasRelation(endpointId models.Id) (bool, error)
}
