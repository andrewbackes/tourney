package service

import (
	"errors"
)

// Common errors that can be returned by the service.
var (
	ErrNotFound  = errors.New("entity not found")
	ErrMalformed = errors.New("missing information needed to create entity")
)
