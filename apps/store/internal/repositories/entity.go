package repositories

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
	"github.com/ilya-mezentsev/micro-dep/store/internal/services/shared"
)

const (
	accountEntitiesQuery   = `SELECT id, author_id, name, description FROM account_linked_entity WHERE account_id = $1`
	entitiesEndpointsQuery = `SELECT * FROM entity_endpoint WHERE entity_id = ANY($1)`
)

type (
	entity struct {
		db        *sqlx.DB
		accountId models.Id
	}

	entityProxy struct {
		Id          string `db:"id"`
		AuthorId    string `db:"author_id"`
		Name        string `db:"name"`
		Description string `db:"description"`
	}
)

func newEntity(db *sqlx.DB, accountId models.Id) entity {
	return entity{db: db, accountId: accountId}
}

func (e entity) Create(model shared.Entity) error {
	//TODO implement me
	panic("implement me")
}

func (e entity) ReadAll() ([]shared.Entity, error) {
	var proxies []entityProxy
	err := e.db.Select(&proxies, accountEntitiesQuery, e.accountId)
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
	//TODO implement me
	panic("implement me")
}

func (e entity) Update(model shared.Entity) (shared.Entity, error) {
	//TODO implement me
	panic("implement me")
}

func (e entity) Delete(id models.Id) error {
	//TODO implement me
	panic("implement me")
}

func (e entity) Exists(name string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (e entity) FetchRelations(entityId models.Id) ([]shared.Endpoint, error) {
	//TODO implement me
	panic("implement me")
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
