package session

import (
	"errors"
	"github.com/ilya-mezentsev/micro-dep/shared/errs"
)

var (
	CredentialsNotFound = errors.Join(errors.New("credentials-not-found"), errs.KeyMissingInStorage)
)
