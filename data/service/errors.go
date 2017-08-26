package service

import (
	"errors"
)

// Common errors that can be returned by the service.
var (
	ErrDoesNotExist = errors.New("entity does not exist")
	ErrMalformed    = errors.New("missing information needed to create entity")
)
