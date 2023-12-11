package errs

import "errors"

var (
	IdMissingInStorage  = errors.New("id-missing-in-storage")
	KeyMissingInStorage = errors.New("key-missing-in-storage")
)
