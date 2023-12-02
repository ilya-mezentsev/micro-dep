package entity

import (
	"errors"

	"github.com/ilya-mezentsev/micro-dep/store/internal/services/shared"
)

var (
	TryingToRemoveEndpointThatIsInUse = errors.Join(errors.New("trying-to-remove-endpoint-that-is-in-use"), shared.Conflict)
	TryingToRemoveEntityThatIsUse     = errors.Join(errors.New("trying-to-remove-endpoint-that-is-in-use"), shared.Conflict)
)
