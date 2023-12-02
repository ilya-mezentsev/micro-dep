package controllers

import (
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
	"github.com/ilya-mezentsev/micro-dep/shared/types/operations"
	"github.com/ilya-mezentsev/micro-dep/store/internal/services"
	servicesShared "github.com/ilya-mezentsev/micro-dep/store/internal/services/shared"
)

type Relation struct {
	controllerMixins[servicesShared.Relation]
}

func NewRelation(servicesFactory func(id models.Id) services.Services) Relation {
	return Relation{
		controllerMixins: controllerMixins[servicesShared.Relation]{
			servicesFactory: servicesFactory,
			serviceFn: func(ss services.Services) operations.CRUD[servicesShared.Relation] {
				return ss.Relation()
			},
		},
	}
}
