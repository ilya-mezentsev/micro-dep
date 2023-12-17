package repositories

import (
	"github.com/jmoiron/sqlx"

	sharedRepositories "github.com/ilya-mezentsev/micro-dep/shared/repositories"
	"github.com/ilya-mezentsev/micro-dep/user/internal/services/shared"
)

const createAuthTokenQuery = `INSERT INTO auth_token VALUES(:author_id, :value, :created_at, :expired_at)`

type (
	Token struct {
		sharedRepositories.AuthToken

		db *sqlx.DB
	}

	authTokenProxy struct {
		AuthorId  string `db:"author_id"`
		Value     string `db:"value"`
		CreatedAt int64  `db:"created_at"`
		ExpiredAt int64  `db:"expired_at"`
	}
)

func newToken(db *sqlx.DB) Token {
	return Token{
		AuthToken: sharedRepositories.NewAuthToken(db),
		db:        db,
	}
}

func (t Token) Create(token shared.AuthToken) error {
	_, err := t.db.NamedExec(createAuthTokenQuery, authTokenProxy{}.fromModel(token))

	return err
}

func (atp authTokenProxy) fromModel(token shared.AuthToken) authTokenProxy {
	return authTokenProxy{
		AuthorId:  string(token.AuthorId),
		Value:     token.Value,
		CreatedAt: token.CreatedAt,
		ExpiredAt: token.ExpiredAt,
	}
}
