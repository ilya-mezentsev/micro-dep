package web

import (
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
	slogGin "github.com/samber/slog-gin"

	"github.com/ilya-mezentsev/micro-dep/shared/types/configs"
	"github.com/ilya-mezentsev/micro-dep/user/internal/services"
	"github.com/ilya-mezentsev/micro-dep/user/internal/transport/web/controllers"
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

	apiGroup := r.Group("/api/user")

	sessionController := controllers.NewSession(services.Session(), webSettings)
	apiGroup.GET("/session", sessionController.Get)
	apiGroup.POST("/session", sessionController.Post)
	apiGroup.DELETE("/session", sessionController.Delete)

	registerController := controllers.NewRegister(services.Register())
	apiGroup.POST("/account", registerController.Register)
	apiGroup.GET("/account/:id", registerController.AccountExists)
	apiGroup.POST("/account/:id", registerController.RegisterForAccount)

	fmt.Printf("Listening port %d\n", webSettings.Port)
	err := r.Run(fmt.Sprintf("%s:%d", webSettings.Domain, webSettings.Port))
	if err != nil {
		panic(err)
	}
}
