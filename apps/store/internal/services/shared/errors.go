package shared

import "errors"

var (
	NotFoundById  = errors.New("not-found-by-id")
	Conflict      = errors.New("conflict")
	AlreadyExists = errors.Join(errors.New("already-exists"), Conflict)
)
