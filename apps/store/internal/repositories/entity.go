package repositories

import (
	"database/sql"
	"errors"

	"github.com/frankenbeanies/uuid4"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/ilya-mezentsev/micro-dep/shared/errs"
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
	"github.com/ilya-mezentsev/micro-dep/store/internal/services/shared"
)

const (
	entityQuery            = `SELECT id, account_id, author_id, name, description FROM entity WHERE account_id = $1 AND id = $2`
	entitiesQuery          = `SELECT id, account_id, author_id, name, description FROM entity WHERE account_id = $1`
	entitiesEndpointsQuery = `SELECT * FROM entity_endpoint WHERE entity_id = ANY($1)`
	endpointsInUseQuery    = `SELECT * FROM entity_endpoint ep WHERE entity_id = $1 AND EXISTS(SELECT 1 FROM dependency WHERE to_id = ep.id)`

	addEntityQuery   = `INSERT INTO entity(id, account_id, author_id, name, description) VALUES (:id, :account_id, :author_id, :name, :description)`
	addEndpointQuery = `INSERT INTO entity_endpoint(id, entity_id, kind, address) VALUES(:id, :entity_id, :kind, :address)`

	updateEntityQuery = `UPDATE entity SET description = $2 WHERE id = $1`
	deleteEntityQuery = `DELETE FROM entity WHERE id = $1`

	entityNameExistsQuery = `SELECT EXISTS(SELECT 1 FROM entity WHERE account_id = $1 AND name = $2)`
)

type (
	entity struct {
		db        *sqlx.DB
		accountId models.Id
	}

	entityProxy struct {
		Id          string `db:"id"`
		AuthorId    string `db:"author_id"`
		AccountId   string `db:"account_id"`
		Name        string `db:"name"`
		Description string `db:"description"`
	}
)

func newEntity(db *sqlx.DB, accountId models.Id) entity {
	return entity{db: db, accountId: accountId}
}

func (e entity) Create(model shared.Entity) (shared.Entity, error) {
	tx, err := e.db.Beginx()
	if err != nil {
		return shared.Entity{}, err
	}
	//goland:noinspection ALL
	defer tx.Rollback()

	model.Id = models.Id(uuid4.New().String())

	_, err = tx.NamedExec(addEntityQuery, entityProxy{AccountId: string(e.accountId)}.fromEntity(model))
	if err != nil {
		return shared.Entity{}, err
	}

	for i := range model.Endpoints {
		model.Endpoints[i].Id = models.Id(uuid4.New().String())
		model.Endpoints[i].EntityId = model.Id
		_, err = tx.NamedExec(addEndpointQuery, endpointProxy{}.fromEndpoint(model.Endpoints[i]))
		if err != nil {
			return shared.Entity{}, err
		}
	}

	return model, tx.Commit()
}

func (e entity) ReadAll() ([]shared.Entity, error) {
	var proxies []entityProxy
	err := e.db.Select(&proxies, entitiesQuery, string(e.accountId))
	if err != nil {
		return nil, err
	}

	var endpointsProxies []endpointProxy
	entitiesIds := e.entitiesIds(proxies)
	err = e.db.Select(&endpointsProxies, entitiesEndpointsQuery, pq.Array(entitiesIds))
	if err != nil {
		return nil, err
	}

	entityId2Endpoints := e.entityId2Endpoints(endpointsProxies, len(proxies))
	result := make([]shared.Entity, len(proxies))
	for i, proxy := range proxies {
		entityModel := proxy.toEntity()
		entityModel.Endpoints = entityId2Endpoints[entityModel.Id]

		result[i] = entityModel
	}

	return result, nil
}

func (e entity) entitiesIds(proxies []entityProxy) []string {
	result := make([]string, len(proxies))
	for i, proxy := range proxies {
		result[i] = proxy.Id
	}

	return result
}

func (e entity) entityId2Endpoints(endpointsProxies []endpointProxy, count int) map[models.Id][]shared.Endpoint {
	result := make(map[models.Id][]shared.Endpoint, count)
	for _, ep := range endpointsProxies {
		endpointModel := ep.toEndpoint()
		result[endpointModel.EntityId] = append(result[endpointModel.EntityId], endpointModel)
	}

	return result
}

func (e entity) ReadOne(id models.Id) (shared.Entity, error) {
	var proxy entityProxy
	err := e.db.Get(&proxy, entityQuery, string(e.accountId), string(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = errs.IdMissingInStorage
		}

		return shared.Entity{}, err
	}

	var endpointProxies []endpointProxy
	err = e.db.Select(&endpointProxies, entitiesEndpointsQuery, pq.Array([]string{proxy.Id}))
	if err != nil {
		return shared.Entity{}, err
	}

	result := proxy.toEntity()
	result.Endpoints = make([]shared.Endpoint, len(endpointProxies))
	for i, ep := range endpointProxies {
		result.Endpoints[i] = ep.toEndpoint()
	}

	return result, nil
}

func (e entity) Update(model shared.Entity) (shared.Entity, error) {
	tx, err := e.db.Beginx()
	if err != nil {
		return shared.Entity{}, err
	}
	//goland:noinspection ALL
	defer tx.Rollback()

	_, err = tx.Exec(updateEntityQuery, string(model.Id), model.Description)
	if err != nil {
		return shared.Entity{}, err
	}

	for _, endpointModel := range model.Endpoints {
		_, err = tx.NamedExec(updateEndpointQuery, endpointProxy{}.fromEndpoint(endpointModel))
		if err != nil {
			return shared.Entity{}, err
		}
	}

	return model, tx.Commit()
}

func (e entity) Delete(id models.Id) error {
	_, err := e.db.Exec(deleteEntityQuery, string(id))

	return err
}

func (e entity) Exists(name string) (bool, error) {
	var result bool
	err := e.db.Get(&result, entityNameExistsQuery, string(e.accountId), name)

	return result, err
}

func (e entity) FetchRelations(entityId models.Id) ([]shared.Endpoint, error) {
	var proxies []endpointProxy
	err := e.db.Select(&proxies, endpointsInUseQuery, string(entityId))
	if err != nil {
		return nil, err
	}

	result := make([]shared.Endpoint, len(proxies))
	for i, proxy := range proxies {
		result[i] = proxy.toEndpoint()
	}

	return result, nil
}

func (ep entityProxy) fromEntity(e shared.Entity) entityProxy {
	return entityProxy{
		Id:          string(e.Id),
		AccountId:   ep.AccountId,
		AuthorId:    string(e.AuthorId),
		Name:        e.Name,
		Description: e.Description,
	}
}

func (ep entityProxy) toEntity() shared.Entity {
	return shared.Entity{
		Id:          models.Id(ep.Id),
		AuthorId:    models.Id(ep.AuthorId),
		Name:        ep.Name,
		Description: ep.Description,
		Endpoints:   nil,
	}
}
