package helpers

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

type (
	assertFn func(*http.Response, *http.Request) error

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

// AssertNoErrors ensures there are no errors in the response
func AssertNoErrors(rsp *http.Response, _ *http.Request) (err error) {
	tmp := StdErrorResponse{}
	return firstErr(decodeBody(rsp, &tmp), tmp)
}

// AssertError ensures there are no errors in the response
func AssertError(expectedError string) assertFn {
	return func(rsp *http.Response, _ *http.Request) (err error) {
		tmp := StdErrorResponse{}
		if err = decodeBody(rsp, &tmp); err != nil {
			return errors.Errorf("Could not decode body: %v", err)
		}

		if tmp.Error.Message == "" {
			return errors.Errorf("No error, expecting: %v", expectedError)
		}

		if expectedError != tmp.Error.Message {
			return errors.Errorf("Expecting error %v, got: %v", expectedError, tmp.Error.Message)
		}

		return nil
	}
}
