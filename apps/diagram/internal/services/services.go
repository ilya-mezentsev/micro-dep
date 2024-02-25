package services

import (
	"log/slog"

	"github.com/ilya-mezentsev/micro-dep/diagram/internal/clients"
	"github.com/ilya-mezentsev/micro-dep/diagram/internal/services/diagram/stateful"
	"github.com/ilya-mezentsev/micro-dep/diagram/internal/services/diagram/stateless"
	"github.com/ilya-mezentsev/micro-dep/diagram/internal/services/draw"
)

type Services struct {
	statefulDiagram  stateful.Service
	statelessDiagram stateless.Service
}

func New(clients clients.Clients, logger *slog.Logger) Services {
	drawService := draw.New(clients.D2())

	return Services{
		statefulDiagram: stateful.New(
			clients.Entities(),
			clients.Relations(),
			drawService,
			logger,
		),

		statelessDiagram: stateless.New(
			drawService,
			logger,
		),
	}
}

func (s Services) StatefulDiagram() stateful.Service {
	return s.statefulDiagram
}

func (s Services) StatelessDiagram() stateless.Service {
	return s.statelessDiagram
}
