package model

import "errors"

var (
	ErrNotFound              = errors.New("not found")
	ErrInvalidArgument       = errors.New("invalid argument")
	ErrUnavailable           = errors.New("service unavailable")
	ErrOrderAlreadyPaid      = errors.New("order has been payd")
	ErrOrderAlreadyCompleted = errors.New("order has been completed")
	ErrOrderAlreadyCancelled = errors.New("order has been cancelled")
	ErrUnexpectedOrderStatus = errors.New("unexpected order status")
	ErrInternal              = errors.New("internal error")
	// ErrProducer         = errors.New("error producing OrderPaid")
)
