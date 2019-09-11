package apitest

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

// Verifier is the assertion interface allowing consumers to inject a custom assertion implementation.
// It also allows failure scenarios to be tested within apitest
type Verifier interface {
	Equal(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) bool
	JSONEq(t *testing.T, expected string, actual string, msgAndArgs ...interface{}) bool
	Fail(t *testing.T, failureMessage string, msgAndArgs ...interface{}) bool
}

// testifyVerifier is a verifier that use https://github.com/stretchr/testify to perform assertions
type testifyVerifier struct{}

// JSONEq asserts that two JSON strings are equivalent
func (a testifyVerifier) JSONEq(t *testing.T, expected string, actual string, msgAndArgs ...interface{}) bool {
	return assert.JSONEq(t, expected, actual, msgAndArgs...)
}

// Equal asserts that two objects are equal
func (a testifyVerifier) Equal(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	return assert.Equal(t, expected, actual, msgAndArgs...)
}

// Fail reports a failure
func (a testifyVerifier) Fail(t *testing.T, failureMessage string, msgAndArgs ...interface{}) bool {
	return assert.Fail(t, failureMessage, msgAndArgs...)
}

func newTestifyVerifier() Verifier {
	return testifyVerifier{}
}

// NoopVerifier is a verifier that does not perform verification
type NoopVerifier struct{}

var _ Verifier = NoopVerifier{}

// Equal does not perform any assertion and always returns true
func (n NoopVerifier) Equal(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	return true
}

// JSONEq does not perform any assertion and always returns true
func (n NoopVerifier) JSONEq(t *testing.T, expected string, actual string, msgAndArgs ...interface{}) bool {
	return true
}

// Fail does not perform any assertion and always returns true
func (n NoopVerifier) Fail(t *testing.T, failureMessage string, msgAndArgs ...interface{}) bool {
	return true
}

// IsSuccess is a convenience function to assert on a range of happy path status codes
var IsSuccess Assert = func(response *http.Response, request *http.Request) error {
	if response.StatusCode >= 200 && response.StatusCode < 400 {
		return nil
	}
	return fmt.Errorf("not success. Status code=%d", response.StatusCode)
}

// IsClientError is a convenience function to assert on a range of client error status codes
var IsClientError Assert = func(response *http.Response, request *http.Request) error {
	if response.StatusCode >= 400 && response.StatusCode < 500 {
		return nil
	}
	return fmt.Errorf("not a client error. Status code=%d", response.StatusCode)
}

// IsServerError is a convenience function to assert on a range of server error status codes
var IsServerError Assert = func(response *http.Response, request *http.Request) error {
	if response.StatusCode >= 500 {
		return nil
	}
	return fmt.Errorf("not a server error. Status code=%d", response.StatusCode)
}
