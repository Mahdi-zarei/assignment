package types

import "errors"

var (
	ErrAlreadyExists     = errors.New("object is already present in the db")
	ErrNotFound          = errors.New("no object found in the db")
	ErrInvalidTime       = errors.New("time is not in valid range")
	ErrGiftIDInvalid     = errors.New("gift id is not valid")
	ErrGifterInvalid     = errors.New("gifter id is not valid")
	ErrGifteeInvalid     = errors.New("giftee id is not valid")
	ErrInvalidResponse   = errors.New("response is not valid")
	ErrInvalidPagination = errors.New("invalid pagination parameters")
	ErrInvalidStatus     = errors.New("invalid status parameter")
)
