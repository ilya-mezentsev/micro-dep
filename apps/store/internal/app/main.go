package app

import (
	"github.com/ilya-mezentsev/micro-dep/store/internal/transport/web"
	"os"

	"github.com/ilya-mezentsev/micro-dep/shared/services/config"
	"github.com/ilya-mezentsev/micro-dep/shared/services/db/connection"
	"github.com/ilya-mezentsev/micro-dep/shared/types/configs"
)

type Config struct {
	DB  configs.DB  `json:"db"`
	Web configs.Web `json:"web"`
}

func Main() {
	settings := config.MustParse[Config](os.Getenv("CONFIG_PATH"))
	db := connection.MustGetConnection(settings.DB)
	servicesFactory := NewServicesFactory(db)

	web.Start(settings.Web, servicesFactory.Services)
}
