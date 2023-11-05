package repositories

import (
	"github.com/jmoiron/sqlx"

	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
	endpointService "github.com/ilya-mezentsev/micro-dep/store/internal/services/endpoint"
	entityService "github.com/ilya-mezentsev/micro-dep/store/internal/services/entity"
	relationService "github.com/ilya-mezentsev/micro-dep/store/internal/services/relation"
)

type Repositories struct {
	entity   entityService.Repo
	endpoint endpointService.Repo
	relation relationService.Repo
}

func New(db *sqlx.DB, accountId models.Id) Repositories {
	return Repositories{
		entity:   newEntity(db, accountId),
		endpoint: newEndpoint(db, accountId),
		relation: newRelation(db, accountId),
	}
}

func (r Repositories) Entity() entityService.Repo {
	return r.entity
}

func (r Repositories) Endpoint() endpointService.Repo {
	return r.endpoint
}

func (r Repositories) Relation() relationService.Repo {
	return r.relation
}
