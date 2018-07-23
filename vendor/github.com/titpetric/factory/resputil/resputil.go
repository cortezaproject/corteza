package resputil

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// Exported error messages
const (
	E_EMPTY_TRACE = "no stack trace available"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

type successMessage struct {
	Success struct {
		Message string `json:"message"`
	} `json:"success"`
}

type errorMessage struct {
	Error struct {
		Message string `json:"message"`
		Trace   string `json:"trace,omitempty"`
	} `json:"error"`
}

// Options struct / configuration parameters
type Options struct {
	Pretty bool // formats JSON output with indentation
	Trace  bool // prints a stack backtrace if exists (pkg/errors)
	Logger func(error)
}

var config Options

// SetConfig updates package options in use
func SetConfig(options Options) {
	config = options
}

// getTrace prints the first available stack trace if any
func getTrace(errs ...error) string {
	for _, err := range errs {
		if err != nil {
			terr, ok := err.(stackTracer)
			if ok {
				return fmt.Sprintf("%+v", terr.StackTrace())
			}
		}
	}
	return E_EMPTY_TRACE
}

// error returns a structured error for API responses
func errorResponse(err ...error) errorMessage {
	response := errorMessage{}
	// add stack trace to the response if available and enabled
	response.Error.Message = "Unknown error"
	if len(err) > 0 {
		if config.Trace {
			response.Error.Trace = getTrace(errors.Cause(err[0]), err[0])
		}
		response.Error.Message = err[0].Error()
	}
	return response
}

// Success returns a structured success message for API responses
func Success(success ...string) successMessage {
	response := successMessage{}
	response.Success.Message = "OK"
	if len(success) > 0 {
		response.Success.Message = success[0]
	}
	return response
}

// OK returns the default Success message
func OK() successMessage {
	return Success()
}

// JSON responds with the first non-nil payload, formats error messages
func JSON(w http.ResponseWriter, responses ...interface{}) {
	respond := func(payload interface{}) {
		var result []byte
		var err error
		encode := func(payload interface{}) ([]byte, error) {
			if config.Pretty {
				return json.MarshalIndent(payload, "", "\t")
			}
			return json.Marshal(payload)
		}
		switch value := payload.(type) {
		case []byte:
			result = value
		case error:
			// main key is "error"
			errWithStack := errors.WithStack(value)
			if config.Logger != nil {
				config.Logger(errWithStack)
			}
			result, err = encode(errorResponse(errWithStack))
		case successMessage:
			// main key is "success"
			result, err = encode(value)
		default:
			// main key is "response"
			result, err = encode(struct {
				Response interface{} `json:"response"`
			}{value})
		}
		if err != nil {
			result, _ = encode(errorResponse(errors.WithStack(err)))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	}

	for _, response := range responses {
		switch value := response.(type) {
		case nil:
			// this will match a nil error
			continue
		case func() ([]byte, error):
			result, err := value()
			JSON(w, err, result)
		case func() (interface{}, error):
			result, err := value()
			JSON(w, err, result)
		case func() error:
			err := value()
			if err == nil {
				continue
			}
			respond(err)
		case *error:
			if *value == nil {
				continue
			}
			respond(*value)
		case error:
			respond(value)
		case string:
			if value == "" {
				continue
			}
			respond(value)
		case bool:
			if !value {
				continue
			}
			respond(value)
		case successMessage:
			respond(value)
		default:
			respond(value)
		}
		// Exit on the first output...
		return
	}
	respond(false)
}
