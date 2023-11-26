package repositories

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/ilya-mezentsev/micro-dep/shared/errs"
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
)

const (
	authorizedAccountIdQuery = `
	SELECT a.account_id
	FROM author a
	INNER JOIN auth_token at ON a.id = at.author_id
	WHERE at.value = $1 AND at.expired_at > $2
	`
)

type AuthToken struct {
	db *sqlx.DB
}

func NewAuthToken(db *sqlx.DB) AuthToken {
	return AuthToken{db: db}
}

// AuthorizedAccountId - returns account id if passed token is valid (exists in DB and is not expired)
// probably this method knows "too much" about authorization
func (a AuthToken) AuthorizedAccountId(token string, authorizedTill time.Time) (models.Id, error) {
	var accountId string
	err := a.db.Get(&accountId, authorizedAccountIdQuery, token, authorizedTill.Unix())
	if errors.Is(err, sql.ErrNoRows) {
		err = errs.IdMissingInStorage
	}

	return models.Id(accountId), err
}
