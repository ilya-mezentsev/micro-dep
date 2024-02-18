package auth

import (
	"time"

	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
)

type (
	TokenReaderRepo interface {
		AuthorizedAccountId(token string, authorizedTill time.Time) (AuthorizedIds, error)
		AccountIdExists(accountId models.Id) (bool, error)
	}

	AuthorizedIds struct {
		AuthorId  models.Id
		AccountId models.Id
	}
)
