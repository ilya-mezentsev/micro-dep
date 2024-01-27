package controllers

import (
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
	"github.com/ilya-mezentsev/micro-dep/store/internal/services"
)

type Entity struct {
	controllerMixins[models.Entity]
}

func NewEntity(servicesFactory func(id models.Id) services.Services) Entity {
	return Entity{
		controllerMixins: controllerMixins[models.Entity]{
			servicesFactory: servicesFactory,
			serviceFn: func(ss services.Services) any {
				return ss.Entity()
			},
		},
	}
}
