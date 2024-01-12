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
	relationsQuery = `
	SELECT d.id, d.from_id, d.to_id
	FROM dependency d
	INNER JOIN entity e ON d.from_id = e.id
	WHERE e.account_id = $1`

	addRelationQuery    = `INSERT INTO dependency(id, from_id, to_id) VALUES(:id, :from_id, :to_id)`
	deleteRelationQuery = `DELETE FROM dependency WHERE id = $1`
)

type (
	relation struct {
		db        *sqlx.DB
		accountId models.Id
	}

	relationProxy struct {
		Id     string `db:"id"`
		FromId string `db:"from_id"`
		ToId   string `db:"to_id"`
	}
)

func newRelation(db *sqlx.DB, accountId models.Id) relation {
	return relation{db: db, accountId: accountId}
}

func (r relation) Create(model shared.Relation) (shared.Relation, error) {
	_, err := r.db.NamedExec(addRelationQuery, relationProxy{}.fromRelation(model))

	return model, err
}

func (r relation) ReadAll() ([]shared.Relation, error) {
	var proxies []relationProxy
	err := r.db.Select(&proxies, relationsQuery, string(r.accountId))
	if err != nil {
		return nil, err
	}

	relations := make([]shared.Relation, len(proxies))
	for i, proxy := range proxies {
		relations[i] = proxy.toRelation()
	}

	return relations, nil
}

func (r relation) ReadOne(_ models.Id) (shared.Relation, error) {
	panic("not implemented")
}

func (r relation) Delete(id models.Id) error {
	_, err := r.db.Exec(deleteRelationQuery, string(id))

	return err
}

func (r relation) PartsExist(model shared.Relation) (bool, bool, error) {
	var entityExists bool
	err := r.db.Get(&entityExists, entityIdExistsQuery, string(model.FromEntityId))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = errs.IdMissingInStorage
		}

		return false, false, err
	}

	var endpointExists bool
	err = r.db.Get(&endpointExists, endpointWithIdExistsQuery, string(model.ToEndpointId))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = errs.IdMissingInStorage
		}

		return false, false, err
	}

	return entityExists, endpointExists, err
}

func (rp relationProxy) fromRelation(r shared.Relation) relationProxy {
	return relationProxy{
		Id:     string(r.Id),
		FromId: string(r.FromEntityId),
		ToId:   string(r.ToEndpointId),
	}
}

func (rp relationProxy) toRelation() shared.Relation {
	return shared.Relation{
		Id:           models.Id(rp.Id),
		FromEntityId: models.Id(rp.FromId),
		ToEndpointId: models.Id(rp.ToId),
	}
}
