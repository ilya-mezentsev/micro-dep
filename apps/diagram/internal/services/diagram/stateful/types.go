package stateful

import (
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
)

type (
	EntitiesFetcher interface {
		// Fetch returns one of: entities, api error message or error
		Fetch(accountId models.Id) ([]models.Entity, error)
	}

	RelationsFetcher interface {
		// Fetch returns one of: relations, api error message or error
		Fetch(accountId models.Id) ([]models.Relation, error)
	}
)
