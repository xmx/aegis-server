package mongodsl

import "fmt"

type DSLError struct {
	Err error
}

func (d *DSLError) Error() string {
	return d.Err.Error()
}

func (d *DSLError) Unwrap() error {
	return d.Err
}

func newDSLErrorf(msg string, args ...any) *DSLError {
	return &DSLError{Err: fmt.Errorf(msg, args...)}
}
