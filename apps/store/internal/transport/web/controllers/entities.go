package controllers

import (
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
	"github.com/ilya-mezentsev/micro-dep/store/internal/services"
	servicesShared "github.com/ilya-mezentsev/micro-dep/store/internal/services/shared"
)

type Entity struct {
	controllerMixins[servicesShared.Entity]
}

func NewEntity(servicesFactory func(id models.Id) services.Services) Entity {
	return Entity{
		controllerMixins: controllerMixins[servicesShared.Entity]{
			servicesFactory: servicesFactory,
			serviceFn: func(ss services.Services) any {
				return ss.Entity()
			},
		},
	}
}
