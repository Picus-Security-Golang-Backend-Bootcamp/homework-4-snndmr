package constants

import "errors"

var (
	ErrBookNotFound       = errors.New("err: not found")
	ErrNegativeAmount     = errors.New("err: amount could not be negative")
	ErrBookOutOfStock     = errors.New("err: out of stock")
	ErrBookAlreadyDeleted = errors.New("err: already deleted")
)
