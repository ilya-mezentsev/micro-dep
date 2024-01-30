package draw

import (
	"bytes"
	"text/template"

	"github.com/ilya-mezentsev/micro-dep/diagram/internal/services/shared/types"
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
)

const tpl = `
{{ range $entity := .Entities }}
	{{ $entity.Name }}: {
		{{ range $endpoint := $entity.Endpoints }}
			{{ $endpoint.Id }}: {{ $endpoint.Address }}
		{{ end }}
	}
{{ end }}


{{ range $relation := .Relations }}
	{{ $relation.FromEntityName }} -> {{ $relation.ToEntityName }}.{{ $relation.ToEndpointId }} : {{ $relation.ToEndpointKind }}
{{ end }}
`

type (
	Service struct {
		d2client D2Client
	}

	templateData struct {
		Entities  []models.Entity
		Relations []relation
	}

	relation struct {
		FromEntityName string
		ToEntityName   string
		ToEndpointId   string
		ToEndpointKind string
	}
)

func New(d2client D2Client) Service {
	return Service{
		d2client: d2client,
	}
}

func (s Service) DrawDiagram(rdp types.RelationsDiagramData) (string, error) {

	t := template.Must(template.New("").Parse(tpl))

	var buffer bytes.Buffer
	err := t.Execute(&buffer, s.prepareTemplateDate(rdp))
	if err != nil {
		return "", err
	}

	return s.d2client.Draw(buffer.String())
}

func (s Service) prepareTemplateDate(rdp types.RelationsDiagramData) templateData {
	entityId2entity := make(map[models.Id]models.Entity, len(rdp.Entities))
	endpointId2endpoint := make(map[models.Id]models.Endpoint)
	endpointId2entity := make(map[models.Id]models.Entity)
	for _, entity := range rdp.Entities {
		entityId2entity[entity.Id] = entity

		for _, endpoint := range entity.Endpoints {
			endpointId2endpoint[endpoint.Id] = endpoint
			endpointId2entity[endpoint.Id] = entity
		}
	}

	relations := make([]relation, 0, len(rdp.Relations))
	for _, r := range rdp.Relations {
		relations = append(relations, relation{
			FromEntityName: entityId2entity[r.FromEntityId].Name,
			ToEntityName:   endpointId2entity[r.ToEndpointId].Name,
			ToEndpointId:   string(r.ToEndpointId),
			ToEndpointKind: endpointId2endpoint[r.ToEndpointId].Kind,
		})
	}

	return templateData{
		Entities:  rdp.Entities,
		Relations: relations,
	}
}
