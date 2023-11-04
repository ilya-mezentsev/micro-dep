package shared

import "github.com/ilya-mezentsev/micro-dep/shared/types/models"

type (
	Entity struct {
		Id          models.Id
		AuthorId    models.Id
		Name        string
		Description string
		Endpoints   []Endpoint
	}

	Endpoint struct {
		Id      models.Id
		Kind    string
		Address string
	}
)
