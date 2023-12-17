package session

import (
	"github.com/ilya-mezentsev/micro-dep/shared/services/auth"
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
	"github.com/ilya-mezentsev/micro-dep/user/internal/services/shared"
)

type (
	AuthorRepo interface {
		ById(authorId models.Id) (shared.Author, error)
		ByCredentials(creds shared.AuthorCreds) (shared.Author, error)
	}

	TokenRepo interface {
		auth.TokenReaderRepo
		Create(token shared.AuthToken) error
	}
)
