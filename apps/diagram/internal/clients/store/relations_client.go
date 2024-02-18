package store

import (
	"time"

	"github.com/ilya-mezentsev/micro-dep/diagram/internal/clients/shared"
	"github.com/ilya-mezentsev/micro-dep/shared/transport/middleware"
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
)

type RelationsFetcher struct {
	relationsAddress string
	timeout          time.Duration
}

func NewRelationsFetcher(relationsAddress string, timeout time.Duration) RelationsFetcher {
	return RelationsFetcher{
		relationsAddress: relationsAddress,
		timeout:          timeout,
	}
}

func (rf RelationsFetcher) Fetch(accountId models.Id) ([]models.Relation, error) {
	return fetch[[]models.Relation](shared.Opts{
		Address: rf.relationsAddress,
		Timeout: rf.timeout,
		Headers: map[string]string{
			middleware.HeaderAccountIdName: string(accountId),
		},
	})
}
