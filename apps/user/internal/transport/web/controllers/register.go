package controllers

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/ilya-mezentsev/micro-dep/shared/transport/shared"
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
	"github.com/ilya-mezentsev/micro-dep/user/internal/services/register"
	servicesModels "github.com/ilya-mezentsev/micro-dep/user/internal/services/shared"
)

type Register struct {
	service register.Service
}

func NewRegister(service register.Service) Register {
	return Register{service: service}
}

func (r Register) AccountExists(context *gin.Context) {
	rb := shared.MakeResponseBuilder(context)

	accountExists, err := r.service.AccountExists(models.Id(context.Param("id")))
	if err != nil {
		rb.InternalError(err)
	} else if accountExists {
		rb.EmptyOk()
	} else {
		rb.EmptyNotFound()
	}
}

func (r Register) Register(context *gin.Context) {
	rb := shared.MakeResponseBuilder(context)

	var creds servicesModels.AuthorCreds
	if err := context.ShouldBindJSON(&creds); err != nil {
		rb.ClientError(err)
		return
	}

	author, err := r.service.Register(creds)
	if err != nil {
		if errors.Is(err, register.UsernameExists) {
			rb.ConflictError(err)
		} else {
			rb.InternalError(err)
		}
	} else {
		rb.Ok(author)
	}
}

func (r Register) RegisterForAccount(context *gin.Context) {
	rb := shared.MakeResponseBuilder(context)

	var creds servicesModels.AuthorCreds
	if err := context.ShouldBindJSON(&creds); err != nil {
		rb.ClientError(err)
		return
	}

	author, err := r.service.RegisterForAccount(
		models.Id(context.Param("id")),
		creds,
	)
	if err != nil {
		if errors.Is(err, register.UsernameExists) {
			rb.ConflictError(err)
		} else if errors.Is(err, register.AccountNotFound) {
			rb.NotFoundError(err)
		} else {
			rb.InternalError(err)
		}
	} else {
		rb.Ok(author)
	}
}
