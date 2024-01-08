package register

import "errors"

var (
	UsernameExists  = errors.New("username-exists")
	AccountNotFound = errors.New("account-not-found")
)
