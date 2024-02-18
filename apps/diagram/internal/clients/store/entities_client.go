package store

import (
	"time"

	"github.com/ilya-mezentsev/micro-dep/diagram/internal/clients/shared"
	"github.com/ilya-mezentsev/micro-dep/shared/transport/middleware"
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
)

type EntitiesFetcher struct {
	entitiesAddress string
	timeout         time.Duration
}

func NewEntitiesFetcher(entitiesAddress string, timeout time.Duration) EntitiesFetcher {
	return EntitiesFetcher{
		entitiesAddress: entitiesAddress,
		timeout:         timeout,
	}
}

func (ef EntitiesFetcher) Fetch(accountId models.Id) ([]models.Entity, error) {
	return fetch[[]models.Entity](shared.Opts{
		Address: ef.entitiesAddress,
		Timeout: ef.timeout,
		Headers: map[string]string{
			middleware.HeaderAccountIdName: string(accountId),
		},
	})
}
