package shared

import "errors"

var (
	ExistsError  = errors.New("entity-exists")
	NotFoundById = errors.New("not-found-by-id")
)
