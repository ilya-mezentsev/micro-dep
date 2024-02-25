package stateless

type (
	Entity struct {
		Name         string     `json:"name"`
		Endpoints    []Endpoint `json:"endpoints"`
		Dependencies []Entity   `json:"dependencies"`
	}

	Endpoint struct {
		Kind    string `json:"kind"`
		Address string `json:"address"`
	}
)
