package errors

import (
	"errors"
	"fmt"
)

type (
	Error struct {
		// serves as base error differentiator and informs error handling process
		// if it makes sense to retry later or with different payload (or not at all)
		kind kind

		// simple error message with basic information
		message string

		// captured call stack
		stack []*frame

		// additional details
		meta meta

		// if error wraps another error
		wrap error
	}

	ErrorHandler func(error) error
)

const (
	stackFrameSkip = 4
)

func err(k kind, m string) *Error {
	return &Error{
		kind:    k,
		message: m,

		// skip 3 steps when collecting frames
		stack: collectStack(stackFrameSkip),
	}
}

// Plain returns plain *Error, without call stack and formatted message
func Plain(k kind, m string, a ...interface{}) *Error {
	return &Error{
		kind:    k,
		message: fmt.Sprintf(m, a...),
	}
}

// New returns *Error
func New(k kind, m string, fn ...mfn) *Error {
	return err(k, m).Apply(fn...)
}

// New returns *Error with formatted message
func Newf(k kind, m string, a ...interface{}) *Error {
	return err(k, fmt.Sprintf(m, a...))
}

// safe to show details of this error
func (Error) Safe() bool { return true }

// Error message
func (e Error) Error() string {
	return e.message
}

// Unwrap wrapped error
func (e Error) Unwrap() error {
	return e.wrap
}

// Wrap error
func (e *Error) Wrap(err error) *Error {
	e.wrap = err
	return e
}

// Recollect stack frames
func (e *Error) Stack(skip int) *Error {
	e.stack = collectStack(skip)
	return e
}

// Returns error's meta
func (e Error) Meta() meta {
	return e.meta
}

// Apply applies modifier functions
func (e *Error) Apply(ffn ...mfn) *Error {
	for _, fn := range ffn {
		fn(e)
	}
	return e
}

// Is provided Is() method for equality checking with errors.Is
func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok {
		return false
	}

	return t.kind == e.kind && t.message == e.message
}

// Unwrap is alias for errors.Unwrap so users can avoid importing both errors packages
//
// This function DOES NOT SUPPRESS errors if they are not wrapped!
func Unwrap(err error) error {
	if err != nil && errors.Unwrap(err) != nil {
		return errors.Unwrap(err)
	}

	return err
}

// Is is alias for errors.Is so users can avoid importing both errors packages
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As is alias for errors.As so users can avoid importing both errors packages
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}
