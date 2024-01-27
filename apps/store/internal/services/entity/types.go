package entity

import (
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
	"github.com/ilya-mezentsev/micro-dep/shared/types/operations"
)

type Service = operations.CRUD[models.Entity]

type Repo interface {
	operations.CRUD[models.Entity]
	Exists(name string) (bool, error)
	// FetchRelations - expected that it returns endpoints to which there are relations
	FetchRelations(entityId models.Id) ([]models.Endpoint, error)
}
