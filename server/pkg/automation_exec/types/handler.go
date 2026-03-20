package types

import (
	"context"
	"errors"
	"fmt"
)

type (
	// Phase is a single closure in a multi-phase step.
	Phase func(context.Context, *ExecRequest) (ExecResponse, error)

	// PhasedStepHandler is an optional upgrade interface for step handlers
	// that need mid-step recoverability. The runtime detects it via type switch.
	PhasedStepHandler interface {
		Phases() []Phase
	}

	// RecoverableError wraps an error to signal that execution may be paused
	// and retried rather than immediately routed to a catch handler.
	RecoverableError struct {
		cause error
	}
)

func NewRecoverableError(msg string) *RecoverableError {
	return &RecoverableError{cause: fmt.Errorf("%s", msg)}
}

func WrapRecoverable(err error) *RecoverableError {
	return &RecoverableError{cause: err}
}

func (e *RecoverableError) Error() string { return e.cause.Error() }
func (e *RecoverableError) Unwrap() error { return e.cause }

// IsRecoverable reports whether err (or any error in its chain) is a RecoverableError.
func IsRecoverable(err error) bool {
	var r *RecoverableError
	return errors.As(err, &r)
}
