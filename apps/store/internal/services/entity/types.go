package entity

import (
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
	"github.com/ilya-mezentsev/micro-dep/shared/types/operations"
	"github.com/ilya-mezentsev/micro-dep/store/internal/services/shared"
)

type Service = operations.CRUD[shared.Entity]

type Repo interface {
	operations.CRUD[shared.Entity]
	Exists(name string) (bool, error)
	// FetchRelations - expected that it returns endpoints to which there are relations
	FetchRelations(entityId models.Id) ([]shared.Endpoint, error)
}
