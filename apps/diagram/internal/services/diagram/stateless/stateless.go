package stateless

import (
	"log/slog"

	"github.com/frankenbeanies/uuid4"

	"github.com/ilya-mezentsev/micro-dep/diagram/internal/services/diagram/shared"
	"github.com/ilya-mezentsev/micro-dep/diagram/internal/services/shared/types"
	"github.com/ilya-mezentsev/micro-dep/shared/errs"
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
)

type Service struct {
	drawService shared.DrawService
	logger      *slog.Logger
	idFactory   func() models.Id
}

func New(
	drawService shared.DrawService,
	logger *slog.Logger,
) Service {

	return Service{
		drawService: drawService,
		logger:      logger,
		idFactory:   defaultIdFactory,
	}
}

// Draw draws diagram by passed slice of Entity
func (s Service) Draw(entities []Entity) (string, error) {
	diagramFilePath, err := s.drawService.DrawDiagram(s.buildRelationsDiagramData(entities))
	if err != nil {
		s.logger.Error(
			"Got an error while drawing diagram",
			slog.Any("error", err),
			slog.Any("entities", entities),
		)

		return "", errs.Unknown
	}

	return diagramFilePath, nil
}

// buildRelationsDiagramData builds relations for drawing diagram
func (s Service) buildRelationsDiagramData(entities []Entity) types.RelationsDiagramData {
	entityName2EntityModel := make(map[string]models.Entity, len(entities))
	relationsModels := make([]models.Relation, 0, len(entities))

	for _, entity := range entities {
		entityModel, ok := entityName2EntityModel[entity.Name]
		if !ok {
			entityModel = s.buildEntityModel(entity)

			entityName2EntityModel[entityModel.Name] = entityModel
		}

		for _, dependency := range entity.Dependencies {
			var endpointsModels []models.Endpoint
			dependencyModel, existingEndpoint := entityName2EntityModel[dependency.Name]
			if existingEndpoint {
				dependencyModel = s.mergeEndpoints(dependencyModel, dependency.Endpoints)
				endpointsModels = s.filterEndpoints(dependencyModel, dependency.Endpoints)
			} else {
				dependencyModel = s.buildEntityModel(dependency)
				endpointsModels = dependencyModel.Endpoints
			}

			entityName2EntityModel[dependency.Name] = dependencyModel

			for _, endpointModel := range endpointsModels {
				relationsModels = append(relationsModels, models.Relation{
					Id:           s.idFactory(),
					FromEntityId: entityModel.Id,
					ToEndpointId: endpointModel.Id,
				})
			}
		}
	}

	entitiesModels := make([]models.Entity, 0, len(entityName2EntityModel))
	for _, entityModel := range entityName2EntityModel {
		entitiesModels = append(entitiesModels, entityModel)
	}

	return types.RelationsDiagramData{
		Entities:  entitiesModels,
		Relations: relationsModels,
	}
}

// buildEntityModel converts Entity to models.Entity.
// Mainly by generating ids
func (s Service) buildEntityModel(entity Entity) models.Entity {
	entityModel := models.Entity{
		Id:        s.idFactory(),
		Name:      entity.Name,
		Endpoints: make([]models.Endpoint, 0, len(entity.Endpoints)),
	}
	for _, endpoint := range entity.Endpoints {
		entityModel.Endpoints = append(entityModel.Endpoints, models.Endpoint{
			Id:       s.idFactory(),
			EntityId: entityModel.Id,
			Kind:     endpoint.Kind,
			Address:  endpoint.Address,
		})
	}

	return entityModel
}

// mergeEndpoints merges endpoints from models.Entity with slice of Endpoint.
// Endpoints identification is performed by Endpoint.Kind and Endpoint.Address
func (s Service) mergeEndpoints(entityModel models.Entity, endpoints []Endpoint) models.Entity {
	kindAndAddress2Endpoint := make(map[string]Endpoint, len(endpoints))
	for _, endpoint := range endpoints {
		kindAndAddress2Endpoint[endpoint.Kind+"+"+endpoint.Address] = endpoint
	}

	kindAndAddress2EndpointModel := make(map[string]models.Endpoint, len(entityModel.Endpoints))
	for _, endpointModel := range entityModel.Endpoints {
		kindAndAddress2EndpointModel[endpointModel.Kind+"+"+endpointModel.Address] = endpointModel
	}

	result := entityModel
	for kindAndAddress, endpoint := range kindAndAddress2Endpoint {
		if _, ok := kindAndAddress2EndpointModel[kindAndAddress]; !ok {
			result.Endpoints = append(result.Endpoints, models.Endpoint{
				Id:       s.idFactory(),
				EntityId: entityModel.Id,
				Kind:     endpoint.Kind,
				Address:  endpoint.Address,
			})
		}
	}

	return result
}

// filterEndpoints filters endpoints from models.Entity
// by comparing kinds and addresses.
func (s Service) filterEndpoints(entityModel models.Entity, endpoints []Endpoint) []models.Endpoint {
	result := make([]models.Endpoint, 0, len(endpoints))
	for _, endpoint := range endpoints {
		for _, endpointModel := range entityModel.Endpoints {
			if endpoint.Kind == endpointModel.Kind && endpoint.Address == endpointModel.Address {
				result = append(result, endpointModel)
				break
			}
		}
	}

	return result
}

func defaultIdFactory() models.Id {
	return models.Id(uuid4.New().String())
}
