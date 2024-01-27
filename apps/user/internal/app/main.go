package app

import (
	"log/slog"
	"os"

	"github.com/ilya-mezentsev/micro-dep/shared/services/config"
	"github.com/ilya-mezentsev/micro-dep/shared/services/db/connection"
	"github.com/ilya-mezentsev/micro-dep/shared/types/configs"
	"github.com/ilya-mezentsev/micro-dep/user/internal/repositories"
	"github.com/ilya-mezentsev/micro-dep/user/internal/services"
	"github.com/ilya-mezentsev/micro-dep/user/internal/transport/web"
)

type Config struct {
	DB  configs.DB  `json:"db"`
	Web configs.Web `json:"web"`
}

func Main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	settings := config.MustParse[Config](os.Getenv("CONFIG_PATH"))
	db := connection.MustGetConnection(settings.DB, logger)
	repos := repositories.New(db)

	ss := services.New(repos, logger)

	web.Start(
		settings.Web,
		ss,
		logger,
	)
}
