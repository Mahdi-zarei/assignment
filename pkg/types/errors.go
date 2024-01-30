package types

import "errors"

var (
	ErrAlreadyExists = errors.New("object is already present in the db")
	ErrNotFound      = errors.New("no object found in the db")
	ErrInvalidTime   = errors.New("time is not in valid range")
)
