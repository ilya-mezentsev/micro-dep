package services

import (
	"log/slog"

	"github.com/ilya-mezentsev/micro-dep/diagram/internal/clients"
	"github.com/ilya-mezentsev/micro-dep/diagram/internal/services/diagram"
	"github.com/ilya-mezentsev/micro-dep/diagram/internal/services/draw"
)

type Services struct {
	diagram diagram.Service
}

func New(clients clients.Clients, logger *slog.Logger) Services {
	return Services{
		diagram: diagram.New(
			clients.Entities(),
			clients.Relations(),
			draw.New(clients.D2()),
			logger,
		),
	}
}

func (s Services) Diagram() diagram.Service {
	return s.diagram
}
