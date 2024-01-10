package relation

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

func (s ServiceImpl) Create(model shared.Relation) (shared.Relation, error) {
	entityExists, endpointExists, err := s.repo.PartsExist(model)
	if err != nil {
		if errors.Is(err, errs.IdMissingInStorage) {
			err = shared.NotFoundById
		} else {
			s.logger.Error(
				"Got an error while checking relation parts existence",
				slog.Any("err", err),
				slog.String("from-entity-id", string(model.FromEntityId)),
				slog.String("to-endpoint-id", string(model.ToEndpointId)),
			)

			err = errs.Unknown
		}

		return shared.Relation{}, err
	} else if !entityExists {
		return shared.Relation{}, TryingToCreateRelationFromMissedEntity
	} else if !endpointExists {
		return shared.Relation{}, TryingToCreateRelationToMissedEndpoint
	}

	model.Id = models.Id(uuid4.New().String())
	relation, err := s.repo.Create(model)
	if err != nil {
		s.logger.Error(
			"Got an error while creating relation",
			slog.Any("err", err),
			slog.String("from-entity-id", string(model.FromEntityId)),
			slog.String("to-endpoint-id", string(model.ToEndpointId)),
		)

		err = errs.Unknown
	}

	return relation, err
}

func (s ServiceImpl) ReadAll() ([]shared.Relation, error) {
	relations, err := s.repo.ReadAll()
	if err != nil {
		s.logger.Error("Got an error while reading all relations", slog.Any("err", err))
		err = errs.Unknown
	}

	return relations, err
}

func (s ServiceImpl) ReadOne(_ models.Id) (shared.Relation, error) {
	panic("not implemented")
}

func (s ServiceImpl) Delete(id models.Id) error {
	err := s.repo.Delete(id)
	if err != nil {
		s.logger.Error(
			"Got an error while deleting relation",
			slog.Any("err", err),
			slog.String("relation-id", string(id)),
		)

		err = errs.Unknown
	}

	return err
}
