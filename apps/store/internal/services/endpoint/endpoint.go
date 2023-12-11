package endpoint

import (
	"errors"

	"github.com/frankenbeanies/uuid4"

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

func (s ServiceImpl) Create(model shared.Endpoint) (shared.Endpoint, error) {
	entityExists, endpointExists, err := s.repo.Exists(model)
	if err != nil {
		if errors.Is(err, errs.IdMissingInStorage) {
			err = shared.NotFoundById
		}

		return shared.Endpoint{}, err
	} else if !entityExists {
		return shared.Endpoint{}, TryingToAddEndpointToMissingEntity
	} else if endpointExists {
		return shared.Endpoint{}, TryingToCreateEndpointThatExists
	}

	model.Id = models.Id(uuid4.New().String())

	return s.repo.Create(model)
}

func (s ServiceImpl) Update(model shared.Endpoint) (shared.Endpoint, error) {
	_, endpointExists, err := s.repo.Exists(model)
	if err != nil {
		if errors.Is(err, errs.IdMissingInStorage) {
			err = shared.NotFoundById
		}

		return shared.Endpoint{}, err
	} else if !endpointExists {
		return shared.Endpoint{}, TryingToUpdateMissingEndpoint
	}

	return s.repo.Update(model)
}

func (s ServiceImpl) Delete(id models.Id) error {
	relationExists, err := s.repo.HasRelation(id)
	if err != nil {
		if errors.Is(err, errs.IdMissingInStorage) {
			err = shared.NotFoundById
		}

		return err
	} else if relationExists {
		return TryingToRemoveEndpointThatHasRelation
	}

	return s.repo.Delete(id)
}
