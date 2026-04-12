package model

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound   = errors.New("not found")
	ErrParseValue = errors.New("error parsing value")
)

type MetadataParseValueError struct {
	Value any
	Kind  error
	err   error
}

func NewMetadataParseValueError(v any, err error) *MetadataParseValueError {
	return &MetadataParseValueError{
		Value: v,
		Kind:  ErrParseValue,
		err:   err,
	}
}

func (e *MetadataParseValueError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%v: %v (value: type=%T %#v)", e.Kind, e.err, e.Value, e.Value)
	}
	return fmt.Sprintf("%v (value: type=%T %#v)", e.Kind, e.Value, e.Value)
}

func (e *MetadataParseValueError) Unwrap() error {
	return e.err
}

func (e *MetadataParseValueError) Is(err error) bool {
	return e.Kind == err
}
