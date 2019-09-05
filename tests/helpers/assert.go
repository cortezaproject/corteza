package helpers

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

type (
	Assert func(*http.Response, *http.Request) error

	StdErrorResponse struct{ Error struct{ Message string } }
)

// decodes response body to given struct
func decodeBody(rsp *http.Response, s interface{}) error {
	if err := json.NewDecoder(rsp.Body).Decode(&s); err != nil {
		return errors.Wrap(err, "could not assert IsAuthorized")
	}

	return nil
}

// Returns first input that could be an error
func firstErr(ee ...interface{}) error {
	for _, e := range ee {
		switch t := e.(type) {
		case error:
			if t != nil {
				return t
			}
		case StdErrorResponse:
			if t.Error.Message != "" {
				return errors.New(t.Error.Message)
			}
		case string:
			if t != "" {
				return errors.New(t)
			}
		}
	}

	return nil
}

// Ensures there are no errors in the response
func NoErrors(rsp *http.Response, _ *http.Request) (err error) {
	tmp := StdErrorResponse{}
	return firstErr(decodeBody(rsp, &tmp), tmp)
}
