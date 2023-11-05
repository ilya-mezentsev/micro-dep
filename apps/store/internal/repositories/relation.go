package repositories

import (
	"github.com/jmoiron/sqlx"

	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
	"github.com/ilya-mezentsev/micro-dep/store/internal/services/shared"
)

type relation struct {
	db        *sqlx.DB
	accountId models.Id
}

func newRelation(db *sqlx.DB, accountId models.Id) relation {
	return relation{db: db, accountId: accountId}
}

func (r relation) Create(model shared.Relation) error {
	//TODO implement me
	panic("implement me")
}

func (r relation) ReadAll() ([]shared.Relation, error) {
	//TODO implement me
	panic("implement me")
}

func (r relation) ReadOne(id models.Id) (shared.Relation, error) {
	//TODO implement me
	panic("implement me")
}

func (r relation) Delete(id models.Id) error {
	//TODO implement me
	panic("implement me")
}

func (r relation) PartsExist(model shared.Relation) (bool, bool, error) {
	//TODO implement me
	panic("implement me")
}
