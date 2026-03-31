package model

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrInvalidArgument  = errors.New("invalid argument")
	ErrUnavailable      = errors.New("service unavailable")
	ErrOrderAlreadyPaid = errors.New("order has been payd")
	ErrInternal         = errors.New("internal error")
)
