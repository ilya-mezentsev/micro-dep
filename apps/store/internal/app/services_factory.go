package app

import (
	"sync"

	"github.com/jmoiron/sqlx"

	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
	"github.com/ilya-mezentsev/micro-dep/store/internal/repositories"
	"github.com/ilya-mezentsev/micro-dep/store/internal/services"
)

type ServicesFactory struct {
	sync.Mutex
	db               *sqlx.DB
	account2services map[models.Id]services.Services
}

func (sf *ServicesFactory) Services(accountId models.Id) services.Services {
	sf.Lock()
	defer sf.Unlock()

	if ss, ok := sf.account2services[accountId]; ok {
		return ss
	}

	repos := repositories.New(sf.db, accountId)
	ss := services.New(repos)
	sf.account2services[accountId] = ss

	return ss
}
