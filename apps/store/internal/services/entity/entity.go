package entity

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

func (s ServiceImpl) Create(model models.Entity) (models.Entity, error) {
	exists, err := s.repo.Exists(model.Name)
	if err != nil {
		s.logger.Error("Got an error while checking entity existence", slog.Any("err", err))

		return models.Entity{}, errs.Unknown
	} else if exists {
		return models.Entity{}, shared.AlreadyExists
	}

	model.Id = models.Id(uuid4.New().String())
	for i := range model.Endpoints {
		model.Endpoints[i].EntityId = model.Id
		model.Endpoints[i].Id = models.Id(uuid4.New().String())
	}

	entity, err := s.repo.Create(model)
	if err != nil {
		s.logger.Error("Got an error while creating entity", slog.Any("err", err))
		err = errs.Unknown
	}

	return entity, err
}

func (s ServiceImpl) ReadAll() ([]models.Entity, error) {
	entities, err := s.repo.ReadAll()
	if err != nil {
		s.logger.Error("Got an error while reading all entities", slog.Any("err", err))
		err = errs.Unknown
	}

	return entities, err
}

func (s ServiceImpl) ReadOne(id models.Id) (models.Entity, error) {
	m, err := s.repo.ReadOne(id)
	if errors.Is(err, errs.IdMissingInStorage) {
		err = shared.NotFoundById
	} else if err != nil {
		s.logger.Error(
			"Got an error while reading entity by id",
			slog.Any("err", err),
			slog.String("entity-id", string(id)),
		)

		err = errs.Unknown
	}

	return m, err
}

func (s ServiceImpl) Update(model models.Entity) (models.Entity, error) {
	endpointsInUse, err := s.repo.FetchRelations(model.Id)
	if err != nil {
		if errors.Is(err, errs.IdMissingInStorage) {
			err = shared.NotFoundById
		} else {
			s.logger.Error(
				"Got an error while fetching entity relations",
				slog.Any("err", err),
				slog.String("entity-id", string(model.Id)),
			)

			err = errs.Unknown
		}

		return models.Entity{}, err
	} else if !s.checkAllEndpointsInUseRemained(model, endpointsInUse) {
		return models.Entity{}, TryingToRemoveEndpointThatIsInUse
	}

	// fixme: should we consider situation when entity is deleted here?
	entity, err := s.repo.Update(model)
	if err != nil {
		s.logger.Error(
			"Got an error while updating entity",
			slog.Any("err", err),
			slog.String("entity-id", string(model.Id)),
		)

		err = errs.Unknown
	}

	return entity, err
}

func (s ServiceImpl) checkAllEndpointsInUseRemained(model models.Entity, endpointsInUse []models.Endpoint) bool {
	newEndpoints := s.makeEndpointsMap(model.Endpoints)
	for _, endpointInUse := range endpointsInUse {
		if _, ok := newEndpoints[endpointInUse.Id]; !ok {
			return false
		}
	}

	return true
}

func (s ServiceImpl) makeEndpointsMap(endpoints []models.Endpoint) map[models.Id]struct{} {
	result := make(map[models.Id]struct{}, len(endpoints))
	for _, endpoint := range endpoints {
		result[endpoint.Id] = struct{}{}
	}

	return result
}

func (s ServiceImpl) Delete(id models.Id) error {
	endpointsInUse, err := s.repo.FetchRelations(id)
	if err != nil {
		if errors.Is(err, errs.IdMissingInStorage) {
			err = shared.NotFoundById
		} else {
			s.logger.Error(
				"Got an error while fetching entity relations",
				slog.Any("err", err),
				slog.String("entity-id", string(id)),
			)

			err = errs.Unknown
		}

		return err
	} else if len(endpointsInUse) > 0 {
		return TryingToRemoveEntityThatIsUse
	}

	err = s.repo.Delete(id)
	if err != nil {
		s.logger.Error(
			"Got an error while deleting entity",
			slog.Any("err", err),
			slog.String("entity-id", string(id)),
		)

		err = errs.Unknown
	}

	return err
}
