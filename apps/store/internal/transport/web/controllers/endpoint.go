package controllers

import (
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
	"github.com/ilya-mezentsev/micro-dep/store/internal/services"
)

type Endpoint struct {
	controllerMixins[models.Endpoint]
}

func NewEndpoint(servicesFactory func(id models.Id) services.Services) Endpoint {
	return Endpoint{
		controllerMixins: controllerMixins[models.Endpoint]{
			servicesFactory: servicesFactory,
			serviceFn: func(ss services.Services) any {
				return ss.Endpoint()
			},
		},
	}
}
