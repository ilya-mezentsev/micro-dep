package app

import (
	"log/slog"
	"os"
	"time"

	"github.com/ilya-mezentsev/micro-dep/diagram/internal/clients"
	"github.com/ilya-mezentsev/micro-dep/diagram/internal/services"
	"github.com/ilya-mezentsev/micro-dep/diagram/internal/services/shared/types"
	"github.com/ilya-mezentsev/micro-dep/diagram/internal/transport/web"
	"github.com/ilya-mezentsev/micro-dep/shared/services/config"
	"github.com/ilya-mezentsev/micro-dep/shared/types/configs"
)

type Config struct {
	Web   configs.Web           `json:"web"`
	Store types.StoreWebConfigs `json:"store"`
}

func Main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	settings := config.MustParse[Config](os.Getenv("CONFIG_PATH"))
	c := clients.New(
		settings.Store.EntitiesAddress,
		settings.Store.RelationsAddress,
		time.Duration(settings.Store.TimeoutMilliseconds)*time.Millisecond,
	)

	s := services.New(c, logger)

	web.Start(
		settings.Web,
		s,
		logger,
	)
}
