package services

import (
	"github.com/ilya-mezentsev/micro-dep/user/internal/repositories"
	"github.com/ilya-mezentsev/micro-dep/user/internal/services/session"
)

type Services struct {
	session session.Service
}

func New(repositories repositories.Repositories) Services {
	return Services{
		session: session.New(
			repositories.Token(),
			repositories.Author(),
		),
	}
}

func (s Services) Session() session.Service {
	return s.session
}
