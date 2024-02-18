package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ilya-mezentsev/micro-dep/shared/services/auth"
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
)

var (
	NoAuthTokenProvided = errors.New("no-auth-token-provided")
	noAccountIdProvided = errors.New("no-account-id-provided")
	invalidToken        = errors.New("invalid-token")
)

type Auth struct {
	service auth.Service
}

func NewAuth(service auth.Service) Auth {
	return Auth{service: service}
}

func (a Auth) ByCookieToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(CookieTokenName)
		if err != nil {
			_ = c.AbortWithError(http.StatusUnauthorized, NoAuthTokenProvided)
			return
		}

		accountId, err := a.service.IsAuthenticatedToken(token)
		if err != nil {
			if errors.Is(err, auth.AccountNotFoundErr) {
				_ = c.AbortWithError(http.StatusUnauthorized, invalidToken)
			} else {
				c.AbortWithStatus(http.StatusInternalServerError)
			}

			return
		}

		c.Set(AccountIdKey, accountId)
		c.Next()
	}
}

func (a Auth) ByHeaderAccountId() gin.HandlerFunc {
	return func(c *gin.Context) {
		headerAccountId := c.GetHeader(HeaderAccountIdName)
		if len(headerAccountId) < 1 {
			_ = c.AbortWithError(http.StatusUnauthorized, noAccountIdProvided)
			return
		}

		accountId, err := a.service.CheckAccountId(models.Id(headerAccountId))
		if err != nil {
			if errors.Is(err, auth.AccountNotFoundErr) {
				_ = c.AbortWithError(http.StatusUnauthorized, invalidToken)
			} else {
				c.AbortWithStatus(http.StatusInternalServerError)
			}

			return
		}

		c.Set(AccountIdKey, accountId)
		c.Next()
	}
}
