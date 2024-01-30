package clients

import (
	"github.com/ilya-mezentsev/micro-dep/diagram/internal/clients/d2"
	"time"

	"github.com/ilya-mezentsev/micro-dep/diagram/internal/clients/store"
	"github.com/ilya-mezentsev/micro-dep/diagram/internal/services/diagram"
	"github.com/ilya-mezentsev/micro-dep/diagram/internal/services/draw"
)

type Clients struct {
	entitiesClient  diagram.EntitiesFetcher
	relationsClient diagram.RelationsFetcher
	d2client        draw.D2Client
}

func New(
	entitiesAddress string,
	relationsAddress string,
	timeout time.Duration,
) Clients {
	return Clients{
		entitiesClient:  store.NewEntitiesFetcher(entitiesAddress, timeout),
		relationsClient: store.NewRelationsFetcher(relationsAddress, timeout),
		d2client:        d2.New(),
	}
}

func (c Clients) Entities() diagram.EntitiesFetcher {
	return c.entitiesClient
}

func (c Clients) Relations() diagram.RelationsFetcher {
	return c.relationsClient
}

func (c Clients) D2() draw.D2Client {
	return c.d2client
}
