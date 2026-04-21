package model

import "errors"

var (
	ErrNotFound        = errors.New("service error: part not found")
	ErrInvalidArgument = errors.New("service error: invalid argument")
	ErrUnexpected      = errors.New("service error: unexpected error")
)
