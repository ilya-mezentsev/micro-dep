package relation

import "errors"

var (
	TryingToCreateRelationFromMissedEntity = errors.New("trying-to-create-relation-from-missed-entity")
	TryingToCreateRelationToMissedEndpoint = errors.New("trying-to-create-relation-to-missed-endpoint")
)
