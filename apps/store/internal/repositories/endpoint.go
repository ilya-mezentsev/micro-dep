package repositories

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"github.com/ilya-mezentsev/micro-dep/shared/errs"
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
	"github.com/ilya-mezentsev/micro-dep/store/internal/services/shared"
)

const (
	updateEndpointQuery      = `UPDATE entity_endpoint SET kind = :kind, address = :address WHERE id = :id`
	deleteEndpointQuery      = `DELETE FROM entity_endpoint WHERE id = $1`
	endpointExistsQuery      = `SELECT EXISTS(SELECT 1 FROM entity_endpoint WHERE id = $1)`
	entityIdExistsQuery      = `SELECT EXISTS(SELECT 1 FROM entity WHERE id = $1)`
	endpointHasRelationQuery = `SELECT EXISTS(SELECT 1 FROM dependency WHERE to_id = $1)`
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

func (e endpoint) Create(model shared.Endpoint) (shared.Endpoint, error) {
	_, err := e.db.NamedExec(addEndpointQuery, endpointProxy{}.fromEndpoint(model))

	return model, err
}

func (e endpoint) Update(model shared.Endpoint) (shared.Endpoint, error) {
	_, err := e.db.NamedExec(updateEndpointQuery, endpointProxy{}.fromEndpoint(model))

	return model, err
}

func (e endpoint) Delete(id models.Id) error {
	_, err := e.db.Exec(deleteEndpointQuery, string(id))

	return err
}

func (e endpoint) Exists(model shared.Endpoint) (bool, bool, error) {
	var entityExists bool
	err := e.db.Get(&entityExists, entityIdExistsQuery, string(model.EntityId))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = errs.IdMissingInStorage
		}

		return false, false, err
	}

	var endpointExists bool
	err = e.db.Get(&endpointExists, endpointExistsQuery, string(model.Id))
	if errors.Is(err, sql.ErrNoRows) {
		err = errs.IdMissingInStorage
	}

	return entityExists, endpointExists, err
}

func (e endpoint) HasRelation(endpointId models.Id) (bool, error) {
	var hasRelation bool
	err := e.db.Get(&hasRelation, endpointHasRelationQuery, string(endpointId))
	if errors.Is(err, sql.ErrNoRows) {
		err = errs.IdMissingInStorage
	}

	return hasRelation, err
}

func (ep endpointProxy) fromEndpoint(e shared.Endpoint) endpointProxy {
	return endpointProxy{
		Id:       string(e.Id),
		EntityId: string(e.EntityId),
		Kind:     e.Kind,
		Address:  e.Address,
	}
}

func (ep endpointProxy) toEndpoint() shared.Endpoint {
	return shared.Endpoint{
		Id:       models.Id(ep.Id),
		EntityId: models.Id(ep.EntityId),
		Kind:     ep.Kind,
		Address:  ep.Address,
	}
}
