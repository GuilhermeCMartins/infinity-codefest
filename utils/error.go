package utils

import "errors"

var (
	ErrBadRequest       = errors.New("Bad request")
	ErrInternalServer   = errors.New("Internal server error")
	ErrNotFound         = errors.New("Not found")
	ErrMethodNotAllowed = errors.New("Method not allowed")
)
