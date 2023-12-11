package repositories

import (
	"database/sql"
	"errors"
	"github.com/ilya-mezentsev/micro-dep/shared/services/auth"
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/ilya-mezentsev/micro-dep/shared/errs"
)

const (
	authorizedAccountIdQuery = `
	SELECT a.id, a.account_id
	FROM author a
	INNER JOIN auth_token at ON a.id = at.author_id
	WHERE at.value = $1 AND at.expired_at > $2
	`
)

type (
	AuthToken struct {
		db *sqlx.DB
	}

	authorizedIds struct {
		AuthorId  string `db:"id"`
		AccountId string `db:"account_id"`
	}
)

func NewAuthToken(db *sqlx.DB) AuthToken {
	return AuthToken{db: db}
}

// AuthorizedAccountId - returns account id if passed token is valid (exists in DB and is not expired)
// probably this method knows "too much" about authorization
func (a AuthToken) AuthorizedAccountId(token string, authorizedTill time.Time) (auth.AuthorizedIds, error) {
	var ids authorizedIds
	err := a.db.Get(&ids, authorizedAccountIdQuery, token, authorizedTill.Unix())
	if errors.Is(err, sql.ErrNoRows) {
		err = errs.IdMissingInStorage
	}

	return auth.AuthorizedIds{
		AuthorId:  models.Id(ids.AuthorId),
		AccountId: models.Id(ids.AccountId),
	}, err
}
