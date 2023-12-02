package endpoint

import (
	"errors"

	"github.com/ilya-mezentsev/micro-dep/store/internal/services/shared"
)

var (
	TryingToAddEndpointToMissingEntity    = errors.Join(errors.New("trying-to-add-endpoint-to-missing-entity"), shared.Conflict)
	TryingToCreateEndpointThatExists      = errors.Join(errors.New("trying-to-create-endpoint-that-exists"), shared.Conflict)
	TryingToUpdateMissingEndpoint         = errors.Join(errors.New("trying-to-update-missing-endpoint"), shared.Conflict)
	TryingToRemoveEndpointThatHasRelation = errors.Join(errors.New("trying-to-remove-endpoint-that-has-relation"), shared.Conflict)
)
