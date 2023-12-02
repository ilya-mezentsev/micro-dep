package controllers

import (
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
	"github.com/ilya-mezentsev/micro-dep/store/internal/services"
	servicesShared "github.com/ilya-mezentsev/micro-dep/store/internal/services/shared"
)

type Endpoint struct {
	controllerMixins[servicesShared.Endpoint]
}

func NewEndpoint(servicesFactory func(id models.Id) services.Services) Endpoint {
	return Endpoint{
		controllerMixins: controllerMixins[servicesShared.Endpoint]{
			servicesFactory: servicesFactory,
			serviceFn: func(ss services.Services) any {
				return ss.Endpoint()
			},
		},
	}
}
