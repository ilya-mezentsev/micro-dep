package entity

import "errors"

var (
	TryingToRemoveEndpointThatIsInUse = errors.New("trying-to-remove-endpoint-that-is-in-use")
	TryingToRemoveEntityThatIsUse     = errors.New("trying-to-remove-entity=that-in-use")
)
