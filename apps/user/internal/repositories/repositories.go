package repositories

import (
	"github.com/jmoiron/sqlx"

	"github.com/ilya-mezentsev/micro-dep/user/internal/services/session"
)

type Repositories struct {
	author session.AuthorRepo
	token  session.TokenRepo
}

func New(db *sqlx.DB) Repositories {
	return Repositories{
		author: newAuthor(db),
		token:  newToken(db),
	}
}

func (r Repositories) Author() session.AuthorRepo {
	return r.author
}

func (r Repositories) Token() session.TokenRepo {
	return r.token
}
