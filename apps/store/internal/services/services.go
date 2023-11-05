package services

import (
	"github.com/ilya-mezentsev/micro-dep/store/internal/repositories"
	"github.com/ilya-mezentsev/micro-dep/store/internal/services/endpoint"
	"github.com/ilya-mezentsev/micro-dep/store/internal/services/entity"
	"github.com/ilya-mezentsev/micro-dep/store/internal/services/relation"
)

type Services struct {
	entity   entity.Service
	endpoint endpoint.Service
	relation relation.Service
}

func New(repositories repositories.Repositories) Services {
	return Services{
		entity:   entity.NewServiceImpl(repositories.Entity()),
		endpoint: endpoint.NewServiceImpl(repositories.Endpoint()),
		relation: relation.NewServiceImpl(repositories.Relation()),
	}
}

func (s Services) Entity() entity.Service {
	return s.entity
}

func (s Services) Endpoint() endpoint.Service {
	return s.endpoint
}

func (s Services) Relation() relation.Service {
	return s.relation
}
