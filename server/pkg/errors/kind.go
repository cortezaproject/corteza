package errors

import (
	"fmt"
	"net/http"
)

type (
	kind uint8
)

const (
	// internal (unspecified) error
	// these errors yield HTTP status 500 and hide
	// details unless in development environment
	KindInternal kind = iota

	// Data validation error
	KindInvalidData

	// Requested data not found
	KindNotFound

	// Stale data submitted
	// Internally data was already updated or is locked in another session
	KindStaleData

	// Data already exists
	KindDuplicateData

	// Access control
	KindUnauthorized

	// Expecting authenticated user
	KindUnauthenticated

	// External system failure
	KindExternal

	// store error
	KindStore

	// object store (file upload/download)
	KindObjStore

	// automation error
	KindAutomation
)

// translates error kind into http status
func (k kind) httpStatus() int {
	switch k {
	case KindInvalidData:
		return http.StatusBadRequest

	case KindNotFound:
		return http.StatusNotFound

	case KindStaleData:
		return http.StatusConflict

	case KindDuplicateData:
		return http.StatusConflict

	case KindUnauthorized:
		return http.StatusUnauthorized

	case KindUnauthenticated:
		return http.StatusForbidden

	default:
		return http.StatusInternalServerError
	}
}

func Internal(m string, aa ...interface{}) *Error {
	return err(KindInternal, fmt.Sprintf(m, aa...))
}

func Store(m string, aa ...interface{}) *Error {
	return err(KindStore, fmt.Sprintf(m, aa...))
}

func ObjStore(m string, aa ...interface{}) *Error {
	return err(KindObjStore, fmt.Sprintf(m, aa...))
}

func InvalidData(m string, aa ...interface{}) *Error {
	return err(KindInvalidData, fmt.Sprintf(m, aa...))
}

func NotFound(m string, aa ...interface{}) *Error {
	return err(KindNotFound, fmt.Sprintf(m, aa...))
}

func StaleData(m string, aa ...interface{}) *Error {
	return err(KindStaleData, fmt.Sprintf(m, aa...))
}

func DuplicateData(m string, aa ...interface{}) *Error {
	return err(KindDuplicateData, fmt.Sprintf(m, aa...))
}

func Unauthorized(m string, aa ...interface{}) *Error {
	return err(KindUnauthorized, fmt.Sprintf(m, aa...))
}

func Unauthenticated(m string, aa ...interface{}) *Error {
	return err(KindUnauthenticated, fmt.Sprintf(m, aa...))
}

func External(m string, aa ...interface{}) *Error {
	return err(KindExternal, fmt.Sprintf(m, aa...))
}

func Automation(m string, aa ...interface{}) *Error {
	return err(KindAutomation, fmt.Sprintf(m, aa...))
}

func IsKind(err error, k kind) bool {
	t, ok := err.(*Error)
	if !ok {
		return false
	}

	return t.kind == k
}

// IsAny returns true if error is of type *Error
func IsAny(err error) bool {
	var is bool
	if err != nil {
		_, is = err.(*Error)
	}

	return is
}

func IsInternal(err error) bool {
	return IsKind(err, KindInternal)
}

func IsStore(err error) bool {
	return IsKind(err, KindStore)
}

func IsObjStore(err error) bool {
	return IsKind(err, KindObjStore)
}

func IsInvalidData(err error) bool {
	return IsKind(err, KindInvalidData)
}

func IsNotFound(err error) bool {
	return IsKind(err, KindNotFound)
}

func IsStaleData(err error) bool {
	return IsKind(err, KindStaleData)
}

func IsDuplicateData(err error) bool {
	return IsKind(err, KindDuplicateData)
}

func IsUnauthorized(err error) bool {
	return IsKind(err, KindUnauthorized)
}

func IsUnauthenticated(err error) bool {
	return IsKind(err, KindUnauthenticated)
}

func IsExternal(err error) bool {
	return IsKind(err, KindExternal)
}

func IsAutomation(err error) bool {
	return IsKind(err, KindAutomation)
}
