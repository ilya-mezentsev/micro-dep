package repositories

import (
	"github.com/jmoiron/sqlx"

	"github.com/ilya-mezentsev/micro-dep/user/internal/services/register"
	"github.com/ilya-mezentsev/micro-dep/user/internal/services/session"
)

type (
	AuthorRepo interface {
		register.AuthorRepo
		session.AuthorRepo
	}

	Repositories struct {
		account register.AccountRepo
		author  AuthorRepo
		token   session.TokenRepo
	}
)

func New(db *sqlx.DB) Repositories {
	return Repositories{
		account: newAccount(db),
		author:  newAuthor(db),
		token:   newToken(db),
	}
}

func (r Repositories) Account() register.AccountRepo {
	return r.account
}

func (r Repositories) Author() AuthorRepo {
	return r.author
}

func (r Repositories) Token() session.TokenRepo {
	return r.token
}
