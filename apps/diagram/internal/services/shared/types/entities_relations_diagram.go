package types

import "github.com/ilya-mezentsev/micro-dep/shared/types/models"

type RelationsDiagramData struct {
	Entities  []models.Entity
	Relations []models.Relation
}
