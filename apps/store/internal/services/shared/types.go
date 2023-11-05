package shared

import "github.com/ilya-mezentsev/micro-dep/shared/types/models"

type (
	Entity struct {
		Id          models.Id  `json:"id"`
		AuthorId    models.Id  `json:"author_id"`
		Name        string     `json:"name"`
		Description string     `json:"description"`
		Endpoints   []Endpoint `json:"endpoints"`
	}

	Endpoint struct {
		Id       models.Id `json:"id"`
		EntityId models.Id `json:"entity_id"`
		Kind     string    `json:"kind"`
		Address  string    `json:"address"`
	}

	Relation struct {
		Id            models.Id
		FromServiceId models.Id
		ToEndpointId  models.Id
	}
)
