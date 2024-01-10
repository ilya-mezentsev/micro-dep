package services

import (
	"log/slog"

	"github.com/ilya-mezentsev/micro-dep/user/internal/repositories"
	"github.com/ilya-mezentsev/micro-dep/user/internal/services/register"
	"github.com/ilya-mezentsev/micro-dep/user/internal/services/session"
)

type Services struct {
	register register.Service
	session  session.Service
}

func New(
	repositories repositories.Repositories,
	logger *slog.Logger,
) Services {

	return Services{
		register: register.New(
			repositories.Account(),
			repositories.Author(),
			logger,
		),

		session: session.New(
			repositories.Token(),
			repositories.Author(),
			logger,
		),
	}
}

func (s Services) Register() register.Service {
	return s.register
}

func (s Services) Session() session.Service {
	return s.session
}
