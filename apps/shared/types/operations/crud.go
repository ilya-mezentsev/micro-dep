package operations

import "github.com/ilya-mezentsev/micro-dep/shared/types/models"

type CRUD[T any] interface {
	Create(model T) error
	ReadAll() ([]T, error)
	ReadOne(id models.Id) (T, error)
	Update(model T) (T, error)
	Delete(id models.Id) error
}
