package relation

import (
	"errors"

	"github.com/ilya-mezentsev/micro-dep/store/internal/services/shared"
)

var (
	TryingToCreateRelationFromMissedEntity = errors.Join(errors.New("trying-to-create-relation-from-missed-entity"), shared.Conflict)
	TryingToCreateRelationToMissedEndpoint = errors.Join(errors.New("trying-to-create-relation-to-missed-endpoint"), shared.Conflict)
)
