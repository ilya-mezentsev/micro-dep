package web

import (
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
	slogGin "github.com/samber/slog-gin"

	"github.com/ilya-mezentsev/micro-dep/diagram/internal/services"
	"github.com/ilya-mezentsev/micro-dep/diagram/internal/transport/web/controllers"
	"github.com/ilya-mezentsev/micro-dep/shared/types/configs"
)

func Start(
	webSettings configs.Web,
	services services.Services,
	logger *slog.Logger,
) {
	r := gin.New()

	r.Use(
		slogGin.New(logger),
		gin.Recovery(),
	)

	apiGroup := r.Group("/api/diagram")

	diagramController := controllers.NewDiagram(services.Diagram())
	apiGroup.GET("/:id", diagramController.Download)

	fmt.Printf("Listening port %d\n", webSettings.Port)
	err := r.Run(fmt.Sprintf("%s:%d", webSettings.Domain, webSettings.Port))
	if err != nil {
		panic(err)
	}
}
