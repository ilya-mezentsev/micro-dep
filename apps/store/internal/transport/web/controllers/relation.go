package controllers

import (
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
	"github.com/ilya-mezentsev/micro-dep/store/internal/services"
)

type Relation struct {
	controllerMixins[models.Relation]
}

func NewRelation(servicesFactory func(id models.Id) services.Services) Relation {
	return Relation{
		controllerMixins: controllerMixins[models.Relation]{
			servicesFactory: servicesFactory,
			serviceFn: func(ss services.Services) any {
				return ss.Relation()
			},
		},
	}
}
