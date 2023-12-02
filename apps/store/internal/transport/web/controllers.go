package web

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/ilya-mezentsev/micro-dep/shared/types/configs"
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
	"github.com/ilya-mezentsev/micro-dep/store/internal/services"
	"github.com/ilya-mezentsev/micro-dep/store/internal/transport/web/controllers"
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
	apiGroup.GET("/relation", relationController.ReadAll)
	apiGroup.POST("/relation", relationController.Create)
	apiGroup.DELETE("/relation/:id", relationController.Delete)

	fmt.Printf("Listening port %d\n", webSettings.Port)
	err := r.Run(fmt.Sprintf("%s:%d", webSettings.Domain, webSettings.Port))
	if err != nil {
		panic(err)
	}
}
