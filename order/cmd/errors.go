package main

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrUnavailable      = errors.New("service unavailable")
	ErrOrderAlreadyPaid = errors.New("order has been payd")
)
