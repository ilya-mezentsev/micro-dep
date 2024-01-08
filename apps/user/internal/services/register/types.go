package register

import (
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
	"github.com/ilya-mezentsev/micro-dep/user/internal/services/shared"
)

type (
	AccountRepo interface {
		Create(account shared.Account) error
		Exists(id models.Id) (bool, error)
	}

	AuthorRepo interface {
		Create(author shared.Author, password string) error
		UsernameExists(username string) (bool, error)
	}
)
