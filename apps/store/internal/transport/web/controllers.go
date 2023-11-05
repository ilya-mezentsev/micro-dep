package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ilya-mezentsev/micro-dep/shared/types/configs"
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
	"github.com/ilya-mezentsev/micro-dep/store/internal/services"
)

type Response struct {
	Message string `json:"message"`
}

func Start(webSettings configs.Web, servicesFactory func(id models.Id) services.Services) {
	r := gin.New()

	r.Use(gin.Recovery())

	apiGroup := r.Group("/api")

	apiGroup.GET("/entities", func(context *gin.Context) {
		accountId := context.GetHeader("X-Account-Id")
		if accountId == "" {
			context.JSON(http.StatusUnauthorized, Response{Message: "no X-Account-Id header"})
			return
		}

		services := servicesFactory(models.Id(accountId))
		entities, err := services.Entity().ReadAll()
		if err != nil {
			context.JSON(http.StatusInternalServerError, Response{Message: err.Error()})
		} else {
			context.JSON(http.StatusOK, entities)
		}
	})

	fmt.Printf("Listening port %d\n", webSettings.Port)
	r.Run(fmt.Sprintf("%s:%d", webSettings.Domain, webSettings.Port))
}
