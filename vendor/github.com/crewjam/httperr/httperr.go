// package httperr implements an error object that speaks HTTP.
package httperr

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// Error represents an error that can be modeled as an
// http status code.
type Error struct {
	StatusCode   int         // the HTTP status code. If not supplied, http.StatusInternalServerError is used.
	Status       string      // the HTTP status text. If not supplied, http.StatusText(http.StatusCode) is used.
	PrivateError error       // an additional error that is not displayed to the user, but may be logged
	Header       http.Header // extra headers to add to the response (optional)
}

func (h Error) Error() string {
	statusCode, statusText := h.StatusCodeAndText()
	return fmt.Sprintf("%d %s", statusCode, statusText)
}

// StatusCodeAndText returns the status code and text of the error
func (h Error) StatusCodeAndText() (int, string) {
	if h.StatusCode == 0 {
		h.StatusCode = http.StatusInternalServerError
	}
	if h.Status == "" {
		h.Status = http.StatusText(h.StatusCode)
	}
	return h.StatusCode, h.Status
}

type statusCodeAndTexter interface {
	StatusCodeAndText() (int, string)
}

// StatusCodeAndText returns the status code and text of the error
func StatusCodeAndText(err error) (int, string) {
	if err == nil {
		return 200, http.StatusText(200)
	}
	err = TranslateError(err)

	if scater, ok := errors.Cause(err).(statusCodeAndTexter); ok {
		return scater.StatusCodeAndText()
	}

	return http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)
}

// WriteResponse writes an error response to w using the specified status code.
func (h Error) WriteResponse(w http.ResponseWriter) {
	if h.StatusCode == 0 {
		h.StatusCode = http.StatusInternalServerError
	}
	if h.Status == "" {
		h.Status = http.StatusText(h.StatusCode)
	}
	for key, values := range h.Header {
		w.Header().Del(key) // overwrite headers already in the response with the ones specified
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	http.Error(w, h.Status, h.StatusCode)
}

// ResponseWriter is an interface for structs that know how to write themselves
// to a response. This interface is implemented by Error.
type ResponseWriter interface {
	WriteResponse(w http.ResponseWriter)
}

// Write writes the specified error to w. If err is a ResponseWriter the
// WriteResponse method is invoked to produce the response. Otherwise a
// generic 500 Internal Server Error is written.
func Write(w http.ResponseWriter, err error) {
	err = TranslateError(err)

	if wr, ok := errors.Cause(err).(ResponseWriter); ok {
		wr.WriteResponse(w)
		return
	}

	wr := Error{
		PrivateError: err,
		StatusCode:   http.StatusInternalServerError,
	}
	wr.WriteResponse(w)
}

// Wrap returns a copy of err with PrivateError set to privateErr.
//
// Example:
//
//  doc := Document{}
//  if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
//     return httperr.Wrap(httperr.BadRequest, err)
//	}
func Wrap(httpErr Error, privateErr error) Error {
	httpErr.PrivateError = privateErr
	return httpErr
}

type onErrorIndexType int

const onErrorIndex onErrorIndexType = iota

// OnError returns a new http.Request that holds a reference to a
// function that will report an error when returned from a request.
func OnError(r *http.Request, f func(err error)) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), onErrorIndex, f))
}

// ReportError reports the error to the function given in
// OnError.
func ReportError(r *http.Request, err error) {
	if v := r.Context().Value(onErrorIndex); v != nil {
		v.(func(error))(err)
	}
}

// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as HTTP handlers.  If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler object that calls f.
type HandlerFunc func(http.ResponseWriter, *http.Request) error

// ServeHTTP calls f(w, r).
func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := f(w, r); err != nil {
		if v := r.Context().Value(onErrorIndex); v != nil {
			v.(func(error))(err)
		}
		Write(w, err)
	}
}

type Public Error

func (h Public) Error() string {
	return Error(h).Error()
}

func (h Public) WriteResponse(w http.ResponseWriter) {
	if h.PrivateError != nil {
		w.Header().Add("X-Error-Message", h.PrivateError.Error())
	}
	for key, values := range h.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	Error(h).WriteResponse(w)
}

func FriendlyBadRequest(err error) error {
	return Public{
		StatusCode:   http.StatusBadRequest,
		PrivateError: err,
	}
}

var (
	BadRequest                   = Error{StatusCode: 400}
	Unauthorized                 = Error{StatusCode: 401}
	PaymentRequired              = Error{StatusCode: 402}
	Forbidden                    = Error{StatusCode: 403}
	NotFound                     = Error{StatusCode: 404}
	MethodNotAllowed             = Error{StatusCode: 405}
	NotAcceptable                = Error{StatusCode: 406}
	ProxyAuthRequired            = Error{StatusCode: 407}
	RequestTimeout               = Error{StatusCode: 408}
	Conflict                     = Error{StatusCode: 409}
	Gone                         = Error{StatusCode: 410}
	LengthRequired               = Error{StatusCode: 411}
	PreconditionFailed           = Error{StatusCode: 412}
	RequestEntityTooLarge        = Error{StatusCode: 413}
	RequestURITooLong            = Error{StatusCode: 414}
	UnsupportedMediaType         = Error{StatusCode: 415}
	RequestedRangeNotSatisfiable = Error{StatusCode: 416}
	ExpectationFailed            = Error{StatusCode: 417}
	Teapot                       = Error{StatusCode: 418}
	TooManyRequests              = Error{StatusCode: 429}
	InternalServerError          = Error{StatusCode: 500}
	NotImplemented               = Error{StatusCode: 501}
	BadGateway                   = Error{StatusCode: 502}
	ServiceUnavailable           = Error{StatusCode: 503}
	GatewayTimeout               = Error{StatusCode: 504}
	HTTPVersionNotSupported      = Error{StatusCode: 505}
)
