package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	assertFn func(*http.Response, *http.Request) error

	StdErrorResponse struct {
		Error struct {
			Message string
		}
	}

	RecordValueErrorSetResponse struct {
		Error struct {
			Message string
			Details []types.RecordValueError
		}
	}
)

// decodes response body to given struct
func DecodeBody(rsp *http.Response, s interface{}) error {
	if err := json.NewDecoder(rsp.Body).Decode(&s); err != nil {
		return fmt.Errorf("could not decode body: %w", err)
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
	return firstErr(DecodeBody(rsp, &tmp), tmp)
}

// Asserts that all expected errors are returned
//
// Compares each error by Kind, Message and Meta.field
//
// Note: This assertion always expects errors!
func AssertRecordValueError(exp ...*types.RecordValueError) assertFn {
	return func(rsp *http.Response, _ *http.Request) (err error) {
		rcvd := RecordValueErrorSetResponse{}
		if err = DecodeBody(rsp, &rcvd); err != nil {
			return err
		}

		if len(rcvd.Error.Details) == 0 {
			return fmt.Errorf("expecting value errors, none received")
		}

	expLoop:
		for _, e := range exp {
			for _, r := range rcvd.Error.Details {
				if e.Kind != r.Kind {
					continue
				}
				if e.Message != r.Message {
					continue
				}
				if e.Meta["field"] != r.Meta["field"] {
					continue
				}

				// found expected error
				continue expLoop
			}

			// did not find expected error
			return fmt.Errorf("did not find expected error %v", e)
		}

		return nil
	}
}

// Dump can be put into Assert()
func Dump(rsp *http.Response, _ *http.Request) (err error) {
	var payload interface{}
	if err = DecodeBody(rsp, &payload); err != nil {
		return err
	}
	return nil
}

// AssertError ensures there are no errors in the response
func AssertError(expectedError string) assertFn {
	return func(rsp *http.Response, _ *http.Request) (err error) {
		tmp := StdErrorResponse{}
		if err = DecodeBody(rsp, &tmp); err != nil {
			return err
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

// AssertBody compares the raw body to the provided string
func AssertBody(expected string) assertFn {
	return func(rsp *http.Response, _ *http.Request) (err error) {
		bb, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			return err
		}

		got := strings.Trim(string(bb), " \n")
		if expected != got {
			return errors.Errorf("Expecting: %v, got: %v", expected, got)
		}
		return nil
	}
}
