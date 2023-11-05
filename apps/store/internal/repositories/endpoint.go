package repositories

import (
	"github.com/jmoiron/sqlx"

	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
	"github.com/ilya-mezentsev/micro-dep/store/internal/services/shared"
)

type (
	endpoint struct {
		db        *sqlx.DB
		accountId models.Id
	}

	endpointProxy struct {
		Id       string `db:"id"`
		EntityId string `db:"entity_id"`
		Kind     string `db:"kind"`
		Address  string `db:"address"`
	}
)

func newEndpoint(db *sqlx.DB, accountId models.Id) endpoint {
	return endpoint{db: db, accountId: accountId}
}

func (e endpoint) Create(model shared.Endpoint) error {
	//TODO implement me
	panic("implement me")
}

func (e endpoint) Update(model shared.Endpoint) (shared.Endpoint, error) {
	//TODO implement me
	panic("implement me")
}

func (e endpoint) Delete(id models.Id) error {
	//TODO implement me
	panic("implement me")
}

func (e endpoint) Exists(model shared.Endpoint) (bool, bool, error) {
	//TODO implement me
	panic("implement me")
}

func (e endpoint) HasRelation(endpointId models.Id) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (ep endpointProxy) toEndpoint() shared.Endpoint {
	return shared.Endpoint{
		Id:       models.Id(ep.Id),
		EntityId: models.Id(ep.EntityId),
		Kind:     ep.Kind,
		Address:  ep.Address,
	}
}
