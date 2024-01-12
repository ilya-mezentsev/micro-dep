package endpoint

import (
	"errors"
	"log/slog"

	"github.com/frankenbeanies/uuid4"

	"github.com/ilya-mezentsev/micro-dep/shared/errs"
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
	"github.com/ilya-mezentsev/micro-dep/store/internal/services/shared"
)

type ServiceImpl struct {
	repo   Repo
	logger *slog.Logger
}

func NewServiceImpl(repo Repo, logger *slog.Logger) ServiceImpl {
	return ServiceImpl{
		repo:   repo,
		logger: logger,
	}
}

func (s ServiceImpl) Create(model shared.Endpoint) (shared.Endpoint, error) {
	entityExists, endpointExists, err := s.repo.Exists(model)
	if err != nil {
		if errors.Is(err, errs.IdMissingInStorage) {
			err = shared.NotFoundById
		} else {
			s.logger.Error(
				"Got an error while checking endpoint entity existence",
				slog.Any("err", err),
				slog.String("entity-id", string(model.EntityId)),
			)

			err = errs.Unknown
		}

		return shared.Endpoint{}, err
	} else if !entityExists {
		// FIXME. Is not errs.IdMissingInStorage error enough?
		return shared.Endpoint{}, TryingToAddEndpointToMissingEntity
	} else if endpointExists {
		return shared.Endpoint{}, TryingToCreateEndpointThatExists
	}

	model.Id = models.Id(uuid4.New().String())
	endpoint, err := s.repo.Create(model)
	if err != nil {
		s.logger.Error(
			"Got an error while creating endpoint",
			slog.Any("err", err),
			slog.String("entity-id", string(model.EntityId)),
		)

		err = errs.Unknown
	}

	return endpoint, err
}

func (s ServiceImpl) Update(model shared.Endpoint) (shared.Endpoint, error) {
	_, endpointExists, err := s.repo.Exists(model)
	if err != nil {
		if errors.Is(err, errs.IdMissingInStorage) {
			err = shared.NotFoundById
		} else {
			s.logger.Error(
				"Got an error while checking endpoint entity existence",
				slog.Any("err", err),
				slog.String("entity-id", string(model.EntityId)),
			)

			err = errs.Unknown
		}

		return shared.Endpoint{}, err
	} else if !endpointExists {
		// FIXME. Is not errs.IdMissingInStorage error enough?
		return shared.Endpoint{}, TryingToUpdateMissingEndpoint
	}

	endpoint, err := s.repo.Update(model)
	if err != nil {
		s.logger.Error(
			"Got an error while updating endpoint",
			slog.Any("err", err),
			slog.String("entity-id", string(model.EntityId)),
			slog.String("endpoint-id", string(model.Id)),
		)

		err = errs.Unknown
	}

	return endpoint, err
}

func (s ServiceImpl) Delete(id models.Id) error {
	relationExists, err := s.repo.HasRelation(id)
	if err != nil {
		if errors.Is(err, errs.IdMissingInStorage) {
			err = shared.NotFoundById
		} else {
			slog.Error(
				"Got an error while checking endpoint relations",
				slog.Any("err", err),
				slog.String("endpoint-id", string(id)),
			)

			err = errs.Unknown
		}

		return err
	} else if relationExists {
		return TryingToRemoveEndpointThatHasRelation
	}

	err = s.repo.Delete(id)
	if err != nil {
		slog.Error(
			"Got an error while deleting endpoint",
			slog.Any("err", err),
			slog.String("endpoint-id", string(id)),
		)

		err = errs.Unknown
	}

	return err
}
