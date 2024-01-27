package web

import (
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
	slogGin "github.com/samber/slog-gin"

	"github.com/ilya-mezentsev/micro-dep/shared/types/configs"
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
	"github.com/ilya-mezentsev/micro-dep/store/internal/services"
	"github.com/ilya-mezentsev/micro-dep/store/internal/transport/web/controllers"
)

func Start(
	webSettings configs.Web,
	servicesFactory func(id models.Id) services.Services,
	authMiddleware gin.HandlerFunc,
	logger *slog.Logger,
) {

	r := gin.New()

	r.Use(
		slogGin.New(logger),
		gin.Recovery(),
	)

	apiGroup := r.Group("/api/dependencies")

	apiGroup.Use(authMiddleware)

	entityController := controllers.NewEntity(servicesFactory)
	apiGroup.GET("/entities", entityController.ReadAll)
	apiGroup.GET("/entity/:id", entityController.ReadOne)
	apiGroup.POST("/entity", entityController.Create)
	apiGroup.PUT("/entity", entityController.Update)
	apiGroup.DELETE("/entity/:id", entityController.Delete)

	endpointController := controllers.NewEndpoint(servicesFactory)
	apiGroup.POST("/endpoint", endpointController.Create)
	apiGroup.PUT("/endpoint", endpointController.Update)
	apiGroup.DELETE("/endpoint/:id", endpointController.Delete)

	relationController := controllers.NewRelation(servicesFactory)
	apiGroup.GET("/relations", relationController.ReadAll)
	apiGroup.POST("/relation", relationController.Create)
	apiGroup.DELETE("/relation/:id", relationController.Delete)

	fmt.Printf("Listening port %d\n", webSettings.Port)
	err := r.Run(fmt.Sprintf("%s:%d", webSettings.Domain, webSettings.Port))
	if err != nil {
		panic(err)
	}
}
