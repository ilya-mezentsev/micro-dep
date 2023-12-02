package controllers

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/ilya-mezentsev/micro-dep/shared/transport/middleware"
	"github.com/ilya-mezentsev/micro-dep/shared/transport/shared"
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
	"github.com/ilya-mezentsev/micro-dep/shared/types/operations"
	"github.com/ilya-mezentsev/micro-dep/store/internal/services"
	servicesShared "github.com/ilya-mezentsev/micro-dep/store/internal/services/shared"
)

var errNoAccountIdInContext = errors.New("unknown-error")

type controllerMixins[T any] struct {
	servicesFactory func(id models.Id) services.Services
	serviceFn       func(ss services.Services) operations.CRUD[T]
}

func (cm controllerMixins[T]) ReadAll(context *gin.Context) {
	ss, rb, err := cm.prepare(context)
	if err != nil {
		return
	}

	responseModels, err := cm.serviceFn(ss).ReadAll()
	if err != nil {
		cm.onError(rb, err)
	} else {
		rb.Ok(responseModels)
	}
}

func (cm controllerMixins[T]) ReadOne(context *gin.Context) {
	ss, rb, err := cm.prepare(context)
	if err != nil {
		return
	}

	entityModel, err := cm.serviceFn(ss).ReadOne(models.Id(context.Param("id")))
	if err != nil {
		cm.onError(rb, err)
	} else {
		rb.Ok(entityModel)
	}
}

func (cm controllerMixins[T]) Create(context *gin.Context) {
	ss, rb, err := cm.prepare(context)
	if err != nil {
		return
	}

	var model T
	if err = context.ShouldBindJSON(&model); err != nil {
		rb.ClientError(err)
		return
	}

	err = cm.serviceFn(ss).Create(model)
	if err != nil {
		cm.onError(rb, err)
	} else {
		rb.EmptyOk()
	}
}

func (cm controllerMixins[T]) Update(context *gin.Context) {
	ss, rb, err := cm.prepare(context)
	if err != nil {
		return
	}

	var model T
	if err = context.ShouldBindJSON(&model); err != nil {
		rb.ClientError(err)
		return
	}

	updatedEntity, err := cm.serviceFn(ss).Update(model)
	if err != nil {
		cm.onError(rb, err)
	} else {
		rb.Ok(updatedEntity)
	}
}

func (cm controllerMixins[T]) Delete(context *gin.Context) {
	ss, rb, err := cm.prepare(context)
	if err != nil {
		return
	}

	err = cm.serviceFn(ss).Delete(models.Id(context.Param("id")))
	if err != nil {
		cm.onError(rb, err)
	} else {
		rb.EmptyOk()
	}
}

// prepare - create structs for request processing;
// attention! returns error to client if no account_id in context
func (cm controllerMixins[T]) prepare(context *gin.Context) (services.Services, shared.ResponseBuilder, error) {
	rb := shared.MakeResponseBuilder(context)
	accountId, exists := context.Get(middleware.AccountIdKey)
	if !exists {
		rb.InternalError(errNoAccountIdInContext)

		return services.Services{}, nil, errNoAccountIdInContext
	}

	return cm.servicesFactory(accountId.(models.Id)), rb, nil
}

func (cm controllerMixins[T]) onError(rb shared.ResponseBuilder, err error) {
	if errors.Is(err, servicesShared.NotFoundById) {
		rb.NotFoundError(err)
	} else if errors.Is(err, servicesShared.Conflict) {
		rb.ConflictError(err)
	} else {
		rb.InternalError(err)
	}
}
