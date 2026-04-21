package model

import "errors"

var (
	ErrBuildQuery  = errors.New("error building query")
	ErrSelectQuery = errors.New("error select query")
	ErrExecQuery   = errors.New("error execute query")
	ErrNotFound    = errors.New("error row not found")
	ErrDuplicate   = errors.New("error duplicated order id")
)
