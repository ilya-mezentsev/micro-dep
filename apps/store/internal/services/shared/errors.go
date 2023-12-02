package shared

import "errors"

var (
	AlreadyExists = errors.New("already-exists")
	NotFoundById  = errors.New("not-found-by-id")
	Conflict      = errors.New("conflict")
)
