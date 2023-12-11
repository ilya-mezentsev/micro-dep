package shared

import "github.com/ilya-mezentsev/micro-dep/shared/types/models"

type (
	Author struct {
		Id           models.Id `json:"id"`
		AccountId    models.Id `json:"account_id"`
		Username     string    `json:"username"`
		RegisteredAt int64     `json:"registered_at"`
	}

	AuthorCreds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	AuthToken struct {
		AuthorId  models.Id
		Value     string
		CreatedAt int64
		ExpiredAt int64
	}

	AuthResult struct {
		Value     string
		ExpiredAt int64
	}
)
