package operations

import "github.com/ilya-mezentsev/micro-dep/shared/types/models"

type (
	Creator[T any] interface {
		Create(model T) error
	}

	Reader[T any] interface {
		ReadAll() ([]T, error)
		ReadOne(id models.Id) (T, error)
	}

	Updater[T any] interface {
		Update(model T) (T, error)
	}

	Deleter interface {
		Delete(id models.Id) error
	}

	CRUD[T any] interface {
		Creator[T]
		Reader[T]
		Updater[T]
		Deleter
	}

	CUD[T any] interface {
		Creator[T]
		Updater[T]
		Deleter
	}

	CRD[T any] interface {
		Creator[T]
		Reader[T]
		Deleter
	}
)
