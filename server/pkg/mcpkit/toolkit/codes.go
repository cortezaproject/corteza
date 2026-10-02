package toolkit

import "fmt"

// The codes a failed tool call carries on the wire. mcpkit turns every handler
// error into an isError result {code, message, next}; these are the values of
// code. A caller branches on code and reads message; next says what to do.
const (
	CodeFailed          = "failed"           // nothing more specific is known
	CodeInvalidArgument = "invalid_argument" // a declared parameter is missing or malformed
	CodeUnknownArgument = "unknown_argument" // a parameter the tool does not declare
	CodeRiskCapped      = "risk_capped"      // the session's risk ceiling refused the call
	CodeNotFound        = "not_found"        // a referenced resource does not exist or is deleted
	CodeForbidden       = "forbidden"        // the caller's roles do not allow it
	CodeUnauthenticated = "unauthenticated"  // no usable identity on the call
	CodeInvalid         = "invalid"          // the service rejected the payload
	CodeConflict        = "conflict"         // a uniqueness or stale-data clash
	CodeTooLarge        = "too_large"        // the result exceeded the size ceiling
)

// Coded is an error that knows its own code and next step. Argument errors
// from this package carry one; a handler may attach one with WithCode.
type Coded interface {
	error
	ToolErrorCode() string
	ToolErrorNext() string
}

type codedError struct {
	code, next string
	err        error
}

func (e *codedError) Error() string         { return e.err.Error() }
func (e *codedError) Unwrap() error         { return e.err }
func (e *codedError) ToolErrorCode() string { return e.code }
func (e *codedError) ToolErrorNext() string { return e.next }

// WithCode attaches a code and a next step to err. next may be empty; the
// transport fills a default for codes that have one.
func WithCode(err error, code, next string) error {
	if err == nil {
		return nil
	}
	return &codedError{code: code, next: next, err: err}
}

// NotFoundf is fmt.Errorf for a resource the caller named and the server does
// not have, carrying CodeNotFound. The message should say what does exist.
func NotFoundf(format string, a ...any) error {
	return WithCode(fmt.Errorf(format, a...), CodeNotFound, "")
}

// Requiredf is fmt.Errorf for a declared argument that is missing or empty,
// carrying CodeInvalidArgument so the transport points at the parameter docs.
func Requiredf(format string, a ...any) error {
	return WithCode(fmt.Errorf(format, a...), CodeInvalidArgument, "")
}
