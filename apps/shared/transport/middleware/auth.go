package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ilya-mezentsev/micro-dep/shared/services/auth"
)

var (
	noTokenInCookie = errors.New("no-token-in-cookie")
	invalidToken    = errors.New("invalid-token")
)

type Auth struct {
	service auth.Service
}

func NewAuth(service auth.Service) Auth {
	return Auth{service: service}
}

func (a Auth) ByCookie() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(CookieName)
		if err != nil {
			_ = c.AbortWithError(http.StatusUnauthorized, noTokenInCookie)
			return
		}

		accountId, err := a.service.IsAuthenticated(token)
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
