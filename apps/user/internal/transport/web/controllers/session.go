package controllers

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/ilya-mezentsev/micro-dep/shared/services/auth"
	sharedMiddleware "github.com/ilya-mezentsev/micro-dep/shared/transport/middleware"
	"github.com/ilya-mezentsev/micro-dep/shared/transport/shared"
	"github.com/ilya-mezentsev/micro-dep/shared/types/configs"
	"github.com/ilya-mezentsev/micro-dep/user/internal/services/session"
	servicesModels "github.com/ilya-mezentsev/micro-dep/user/internal/services/shared"
)

type Session struct {
	service session.Service
	config  configs.Web
}

func NewSession(service session.Service, config configs.Web) Session {
	return Session{
		service: service,
		config:  config,
	}
}

func (s Session) Get(context *gin.Context) {
	rb := shared.MakeResponseBuilder(context)

	token, err := context.Cookie(sharedMiddleware.CookieName)
	if err != nil {
		rb.UnauthorizedError(sharedMiddleware.NoTokenInCookie)
		return
	}

	author, err := s.service.AuthorizedByToken(token)
	if err != nil {
		if errors.Is(err, auth.AccountNotFoundErr) {
			rb.NotFoundError(err)
		} else {
			rb.InternalError(err)
		}

		return
	}

	rb.Ok(author)
}

func (s Session) Post(context *gin.Context) {
	rb := shared.MakeResponseBuilder(context)

	var creds servicesModels.AuthorCreds
	if err := context.ShouldBindJSON(&creds); err != nil {
		rb.ClientError(err)
		return
	}

	author, authResult, err := s.service.AuthorizeByCredentials(creds)
	if err != nil {
		if errors.Is(err, session.CredentialsNotFound) {
			rb.NotFoundError(err)
		} else {
			rb.InternalError(err)
		}

		return
	}

	context.SetCookie(
		sharedMiddleware.CookieName,
		authResult.Value,
		int(authResult.ExpiredAt),
		"/",
		s.config.Domain,
		s.config.SecureCookie,
		s.config.HttpOnly,
	)

	rb.Ok(author)
}

func (s Session) Delete(context *gin.Context) {
	rb := shared.MakeResponseBuilder(context)

	context.SetCookie(
		sharedMiddleware.CookieName,
		"",
		0,
		"/",
		s.config.Domain,
		s.config.SecureCookie,
		s.config.HttpOnly,
	)

	rb.EmptyOk()
}
