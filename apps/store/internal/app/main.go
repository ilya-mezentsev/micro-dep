package app

import (
	"log/slog"
	"os"

	"github.com/ilya-mezentsev/micro-dep/shared/repositories"
	"github.com/ilya-mezentsev/micro-dep/shared/services/auth"
	"github.com/ilya-mezentsev/micro-dep/shared/services/config"
	"github.com/ilya-mezentsev/micro-dep/shared/services/db/connection"
	"github.com/ilya-mezentsev/micro-dep/shared/transport/middleware"
	"github.com/ilya-mezentsev/micro-dep/shared/types/configs"
	"github.com/ilya-mezentsev/micro-dep/store/internal/transport/web"
)

type Config struct {
	DB  configs.DB  `json:"db"`
	Web configs.Web `json:"web"`
}

func Main() {
	settings := config.MustParse[Config](os.Getenv("CONFIG_PATH"))
	db := connection.MustGetConnection(settings.DB)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	servicesFactory := NewServicesFactory(db, logger)

	authService := auth.NewService(repositories.NewAuthToken(db), logger)
	authMiddleware := middleware.NewAuth(authService)

	web.Start(
		settings.Web,
		servicesFactory.Services,
		authMiddleware.ByCookie(),
		logger,
	)
}
