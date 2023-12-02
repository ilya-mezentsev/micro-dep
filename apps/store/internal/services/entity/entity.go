package entity

import (
	"errors"

	"github.com/ilya-mezentsev/micro-dep/shared/errs"
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
	"github.com/ilya-mezentsev/micro-dep/store/internal/services/shared"
)

type ServiceImpl struct {
	repo Repo
}

func NewServiceImpl(repo Repo) ServiceImpl {
	return ServiceImpl{repo: repo}
}

func (s ServiceImpl) Create(model shared.Entity) error {
	exists, err := s.repo.Exists(model.Name)
	if err != nil {
		return err
	} else if exists {
		return shared.AlreadyExists
	}

	return s.repo.Create(model)
}

func (s ServiceImpl) ReadAll() ([]shared.Entity, error) {
	return s.repo.ReadAll()
}

func (s ServiceImpl) ReadOne(id models.Id) (shared.Entity, error) {
	m, err := s.repo.ReadOne(id)
	if errors.Is(err, errs.IdMissingInStorage) {
		err = shared.NotFoundById
	}

	return m, err
}

func (s ServiceImpl) Update(model shared.Entity) (shared.Entity, error) {
	endpointsInUse, err := s.repo.FetchRelations(model.Id)
	if err != nil {
		if errors.Is(err, errs.IdMissingInStorage) {
			err = shared.NotFoundById
		}

		return shared.Entity{}, err
	} else if !s.checkAllEndpointsInUseRemained(model, endpointsInUse) {
		return shared.Entity{}, TryingToRemoveEndpointThatIsInUse
	}

	// fixme: should we consider situation when entity is deleted here?
	return s.repo.Update(model)
}

func (s ServiceImpl) checkAllEndpointsInUseRemained(model shared.Entity, endpointsInUse []shared.Endpoint) bool {
	newEndpoints := s.makeEndpointsMap(model.Endpoints)
	for _, endpointInUse := range endpointsInUse {
		if _, ok := newEndpoints[endpointInUse.Id]; !ok {
			return false
		}
	}

	return true
}

func (s ServiceImpl) makeEndpointsMap(endpoints []shared.Endpoint) map[models.Id]struct{} {
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
		}

		return err
	} else if len(endpointsInUse) > 0 {
		return TryingToRemoveEntityThatIsUse
	}

	return s.repo.Delete(id)
}
