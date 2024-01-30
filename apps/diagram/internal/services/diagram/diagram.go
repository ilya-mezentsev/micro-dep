package diagram

import (
	"errors"
	"log/slog"
	"sync"

	"github.com/ilya-mezentsev/micro-dep/diagram/internal/services/shared/types"
	"github.com/ilya-mezentsev/micro-dep/shared/errs"
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
)

type Service struct {
	entitiesFetcher  EntitiesFetcher
	relationsFetcher RelationsFetcher
	drawService      DrawService
	logger           *slog.Logger
}

func New(
	entitiesFetcher EntitiesFetcher,
	relationsFetcher RelationsFetcher,
	drawService DrawService,
	logger *slog.Logger,
) Service {

	return Service{
		entitiesFetcher:  entitiesFetcher,
		relationsFetcher: relationsFetcher,
		drawService:      drawService,
		logger:           logger,
	}
}

func (s Service) Draw(accountId models.Id) (string, error) {
	var (
		entities    []models.Entity
		entitiesErr error

		relations    []models.Relation
		relationsErr error

		wg sync.WaitGroup
	)

	wg.Add(1)
	go func() {
		defer wg.Done()

		entities, entitiesErr = s.entitiesFetcher.Fetch(accountId)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		relations, relationsErr = s.relationsFetcher.Fetch(accountId)
	}()

	wg.Wait()

	if entitiesErr != nil || relationsErr != nil {
		s.logger.Error(
			"Got an error while fetching entities",
			slog.Any("entities-error", entitiesErr),
			slog.Any("relations-error", relationsErr),
			slog.String("account-id", string(accountId)),
		)

		return "", errors.Join(entitiesErr, relationsErr)
	}

	diagramFilePath, err := s.drawService.DrawDiagram(types.RelationsDiagramData{
		Entities:  entities,
		Relations: relations,
	})
	if err != nil {
		s.logger.Error(
			"Got an error while drawing diagram",
			slog.Any("entities-error", entitiesErr),
			slog.Any("relations-error", relationsErr),
			slog.String("account-id", string(accountId)),
		)

		return "", errs.Unknown
	}

	return diagramFilePath, nil
}
