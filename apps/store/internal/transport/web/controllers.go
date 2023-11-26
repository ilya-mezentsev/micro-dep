package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ilya-mezentsev/micro-dep/shared/transport/middleware"
	"github.com/ilya-mezentsev/micro-dep/shared/types/configs"
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
	"github.com/ilya-mezentsev/micro-dep/store/internal/services"
)

type Response struct {
	Message string `json:"message"`
}

func Start(
	webSettings configs.Web,
	servicesFactory func(id models.Id) services.Services,
	cookieAuthMiddleware gin.HandlerFunc,
) {

	r := gin.New()

	r.Use(gin.Recovery())

	apiGroup := r.Group("/api")

	apiGroup.Use(cookieAuthMiddleware)

	apiGroup.GET("/entities", func(context *gin.Context) {
		accountId, exists := context.Get(middleware.AccountIdKey)
		if !exists {
			context.JSON(http.StatusInternalServerError, Response{Message: "internal error"})
			return
		}

		accountServices := servicesFactory(accountId.(models.Id))
		entities, err := accountServices.Entity().ReadAll()
		if err != nil {
			context.JSON(http.StatusInternalServerError, Response{Message: err.Error()})
		} else {
			context.JSON(http.StatusOK, entities)
		}
	})

	fmt.Printf("Listening port %d\n", webSettings.Port)
	r.Run(fmt.Sprintf("%s:%d", webSettings.Domain, webSettings.Port))
}
