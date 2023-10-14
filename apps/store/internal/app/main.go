package app

import (
	"github.com/ilya-mezentsev/micro-dep/shared/services/db/connection"
	"github.com/ilya-mezentsev/micro-dep/shared/types/configs"
)

func Main() {
	_ = connection.MustGetConnection(configs.DB{
		Host:     "localhost",
		Port:     5432,
		User:     "user",
		Password: "password",
		DBName:   "dep",
		Connection: struct {
			RetryCount   int `json:"retry_count"`
			RetryTimeout int `json:"retry_timeout"`
		}{
			RetryCount:   2,
			RetryTimeout: 5,
		},
	})
}
