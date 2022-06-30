package errs

import "errors"

var (
	ErrNotFound               = errors.New("not found")
	ErrAlreadyExists          = errors.New("already exists")
	ErrAlreadyUploadThisUser  = errors.New("already upload of this user")
	ErrAlreadyUploadOtherUser = errors.New("already upload of other user")
	ErrBadLoginOrPass         = errors.New("bad login or pass")
)
