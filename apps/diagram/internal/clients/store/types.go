package store

type (
	OkResponse[T any] struct {
		Data T `json:"data"`
	}

	ErrorResponse struct {
		Error string `json:"error"`
	}
)
