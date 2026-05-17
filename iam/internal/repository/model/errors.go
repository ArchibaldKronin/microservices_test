package model

import "errors"

var (
	ErrNotFound    = errors.New("error row not found")
	ErrSelectQuery = errors.New("error select query")
	ErrConverter   = errors.New("error convert data")
	ErrBuildQuery  = errors.New("error building query")
	ErrDuplicate   = errors.New("error duplicated user id")
	ErrExecQuery   = errors.New("error execute query")
)
