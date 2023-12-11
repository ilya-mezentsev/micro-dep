package relation

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

func (s ServiceImpl) Create(model shared.Relation) (shared.Relation, error) {
	entityExists, endpointExists, err := s.repo.PartsExist(model)
	if err != nil {
		if errors.Is(err, errs.IdMissingInStorage) {
			err = shared.NotFoundById
		}

		return shared.Relation{}, err
	} else if !entityExists {
		return shared.Relation{}, TryingToCreateRelationFromMissedEntity
	} else if !endpointExists {
		return shared.Relation{}, TryingToCreateRelationToMissedEndpoint
	}

	model.Id = models.Id(uuid4.New().String())

	return s.repo.Create(model)
}

func (s ServiceImpl) ReadAll() ([]shared.Relation, error) {
	return s.repo.ReadAll()
}

func (s ServiceImpl) ReadOne(_ models.Id) (shared.Relation, error) {
	panic("not implemented")
}

func (s ServiceImpl) Delete(id models.Id) error {
	return s.repo.Delete(id)
}
