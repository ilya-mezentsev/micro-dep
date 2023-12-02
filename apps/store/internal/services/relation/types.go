package relation

import (
	"github.com/ilya-mezentsev/micro-dep/shared/types/operations"
	"github.com/ilya-mezentsev/micro-dep/store/internal/services/shared"
)

type Service = operations.CRUD[shared.Relation]

type Repo interface {
	operations.CRD[shared.Relation]
	// PartsExist - expected to return two flags: entity existence and endpoint existence
	PartsExist(model shared.Relation) (bool, bool, error)
}
