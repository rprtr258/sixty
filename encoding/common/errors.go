package common

import (
	"fmt"
	"reflect"
)

// MismatchError - value does not match its schema, Reason describes why.
type MismatchError struct {
	Reason string
}

func (e MismatchError) Error() string { return e.Reason }

// PathError adds a location (field name or element index) to a nested error.
type PathError struct {
	Path string
	Err  error
}

func (e PathError) Error() string { return e.Path + ": " + e.Err.Error() }

func (e PathError) Unwrap() error { return e.Err }

func kindMismatch(expected string, got reflect.Kind) MismatchError {
	return MismatchError{Reason: fmt.Sprintf("expected %s, got %s", expected, got)}
}

func valueMismatch(expected string, got any) MismatchError {
	return MismatchError{Reason: fmt.Sprintf("expected %s value, got %T", expected, got)}
}
