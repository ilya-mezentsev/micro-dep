package relation

import (
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
	"github.com/ilya-mezentsev/micro-dep/shared/types/operations"
)

type Service = operations.CRD[models.Relation]

type Repo interface {
	operations.CRD[models.Relation]
	// PartsExist - expected to return two flags: entity existence and endpoint existence
	PartsExist(model models.Relation) (bool, bool, error)
}
