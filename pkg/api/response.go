package api

// This code is modified version from
// https://github.com/titpetric/factory/tree/master/resputil
//
// Parts of the code are rewritten to allow greater flexibility
// but the general logic stays the same for now

import (
	"encoding/json"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"net/http"
)

type (
	successWrap struct {
		Success struct {
			Message string `json:"message"`
		} `json:"success"`
	}

	CallFrame struct {
		Function string `json:"function"`
		File     string `json:"file"`
		Line     int    `json:"line"`
	}

	ErrorPayload struct {
		Message   string                 `json:"message"`
		Context   map[string]interface{} `json:"context,omitempty"`
		Callstack []*CallFrame           `json:"callstack,omitempty"`
	}
)

// Success returns a structured success message for API responses
func Success(success ...string) *successWrap {
	response := &successWrap{}
	response.Success.Message = "OK"
	if len(success) > 0 {
		response.Success.Message = success[0]
	}
	return response
}

// OK returns the default Success message
func OK() *successWrap {
	return Success()
}

// Writes response, according to type
//
// Primarily this function encodes given payload (directly or indirectly) as compact JSON
//
// In some specific scenarios, when:
//  - debug mode is enabled,
//  - no explicit accept header with /json mime-type is sent
// and,
//
// if payload is an error:
//    it outputs formatted error with extended info
//
// if payload is non-error:
//    it outputs formatted and indented JSON
//
func encode(w http.ResponseWriter, r *http.Request, payload interface{}) {
	var (
		err error
		enc = json.NewEncoder(w)
	)

	switch c := payload.(type) {

	case error:
		err = c

	case *successWrap:
		// main key is "success"
		w.Header().Add("Content-Type", "application/json")
		if err = enc.Encode(c); err != nil {
			err = fmt.Errorf("failed to encode response: %w", err)
		}

	default:
		w.Header().Add("Content-Type", "application/json")
		// main key is "response"
		aux := struct {
			Response interface{} `json:"response"`
		}{c}
		if err = enc.Encode(aux); err != nil {
			err = fmt.Errorf("failed to encode response: %w", err)
		}
	}

	if err != nil {
		if err, is := err.(*errors.Error); is {
			// trim out the base stack we don't care about...
			_ = err.Apply(errors.StackTrimAtFn("http.HandlerFunc.ServeHTTP"))
		}

		errors.ServeHTTP(w, r, err, !DebugFromContext(r.Context()))
		return
	}
}

// Send handles first non-nil (and non-empty) payload and encodes it or it's results (when fn)
//
// See encode() for details
func Send(w http.ResponseWriter, r *http.Request, rr ...interface{}) {
	for _, rsp := range rr {
		switch c := rsp.(type) {
		case nil:
			// this will match a nil error
			continue

		case *successWrap:
			encode(w, r, c)

		case func(http.ResponseWriter, *http.Request):
			c(w, r)

		case func() ([]byte, error):
			result, err := c()
			Send(w, r, err, result)

		case func() (interface{}, error):
			result, err := c()
			Send(w, r, err, result)

		case func() error:
			err := c()
			if err == nil {
				continue
			}
			encode(w, r, err)

		case error:
			encode(w, r, c)

		case []byte:
			if len(c) == 0 {
				continue
			}
			if _, err := w.Write(c); err != nil {
				panic(err)
			}
			return

		case string:
			if c == "" {
				continue
			}
			encode(w, r, c)

		case bool:
			if !c {
				continue
			}
			encode(w, r, c)

		default:
			encode(w, r, c)
		}

		// Exit on the first output...
		return
	}

	encode(w, r, false)
}
