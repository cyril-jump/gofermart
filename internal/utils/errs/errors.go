package errs

import "errors"

var (
	ErrNotFound       = errors.New("not found")
	ErrAlreadyExists  = errors.New("already exists")
	ErrBadLoginOrPass = errors.New("bad login or pass")
)
