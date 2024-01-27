package models

type (
	Entity struct {
		Id          Id         `json:"id"`
		AuthorId    Id         `json:"author_id"`
		Name        string     `json:"name"`
		Description string     `json:"description"`
		Endpoints   []Endpoint `json:"endpoints"`
	}

	Endpoint struct {
		Id       Id     `json:"id"`
		EntityId Id     `json:"entity_id"`
		Kind     string `json:"kind"`
		Address  string `json:"address"`
	}

	Relation struct {
		Id           Id `json:"id"`
		FromEntityId Id `json:"from_entity_id"`
		ToEndpointId Id `json:"to_endpoint_id"`
	}
)
