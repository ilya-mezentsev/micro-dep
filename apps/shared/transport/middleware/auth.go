package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ilya-mezentsev/micro-dep/shared/services/auth"
)

var (
	NoAuthTokenProvided = errors.New("no-auth-token-provided")
	invalidToken        = errors.New("invalid-token")
)

type Auth struct {
	service auth.Service
}

func NewAuth(service auth.Service) Auth {
	return Auth{service: service}
}

func (a Auth) ByToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, ok := a.token(c)
		if !ok {
			_ = c.AbortWithError(http.StatusUnauthorized, NoAuthTokenProvided)
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

func (a Auth) token(c *gin.Context) (token string, isTokenFound bool) {
	token, err := c.Cookie(TokenName)
	if err != nil {
		token = c.GetHeader(TokenName)
		isTokenFound = len(token) > 0
	} else {
		isTokenFound = true
	}

	return token, isTokenFound
}
