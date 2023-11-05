package repositories

import (
	"github.com/jmoiron/sqlx"

	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
	"github.com/ilya-mezentsev/micro-dep/store/internal/services/shared"
)

type entity struct {
	db        *sqlx.DB
	accountId models.Id
}

func newEntity(db *sqlx.DB, accountId models.Id) entity {
	return entity{db: db, accountId: accountId}
}

func (e entity) Create(model shared.Entity) error {
	//TODO implement me
	panic("implement me")
}

func (e entity) ReadAll() ([]shared.Entity, error) {
	//TODO implement me
	panic("implement me")
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
