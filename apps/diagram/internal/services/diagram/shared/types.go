package shared

import "github.com/ilya-mezentsev/micro-dep/diagram/internal/services/shared/types"

type DrawService interface {
	// DrawDiagram returns path to diagram file or error
	DrawDiagram(rdp types.RelationsDiagramData) (string, error)
}
