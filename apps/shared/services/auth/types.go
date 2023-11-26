package auth

import (
	"time"

	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
)

type TokenRepo interface {
	AuthorizedAccountId(token string, authorizedTill time.Time) (models.Id, error)
}
