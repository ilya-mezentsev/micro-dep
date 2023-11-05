package endpoint

import "errors"

var (
	TryingToAddEndpointToMissingEntity    = errors.New("trying-to-add-endpoint-to-missing-entity")
	TryingToCreateEndpointThatExists      = errors.New("trying-to-create-endpoint-that-exists")
	TryingToUpdateMissingEndpoint         = errors.New("trying-to-update-missing-endpoint")
	TryingToRemoveEndpointThatHasRelation = errors.New("trying-to-remove-endpoint-that-has-relation")
)
