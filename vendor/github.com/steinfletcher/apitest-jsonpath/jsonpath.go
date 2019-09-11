package jsonpath

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	regex "regexp"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/steinfletcher/apitest"
)

// Contains is a convenience function to assert that a jsonpath expression extracts a value in an array
func Contains(expression string, expected interface{}) apitest.Assert {
	return func(res *http.Response, req *http.Request) error {
		value, err := jsonPath(res.Body, expression)
		if err != nil {
			return err
		}

		ok, found := includesElement(value, expected)
		if !ok {
			return errors.New(fmt.Sprintf("\"%s\" could not be applied builtin len()", expected))
		}
		if !found {
			return errors.New(fmt.Sprintf("\"%s\" does not contain \"%s\"", expected, value))
		}
		return nil
	}
}

// Equal is a convenience function to assert that a jsonpath expression extracts a value
func Equal(expression string, expected interface{}) apitest.Assert {
	return func(res *http.Response, req *http.Request) error {
		value, err := jsonPath(res.Body, expression)
		if err != nil {
			return err
		}

		if !objectsAreEqual(value, expected) {
			return errors.New(fmt.Sprintf("\"%s\" not equal to \"%s\"", value, expected))
		}
		return nil
	}
}

func Len(expression string, expectedLength int) apitest.Assert {
	return func(res *http.Response, req *http.Request) error {
		value, err := jsonPath(res.Body, expression)
		if err != nil {
			return err
		}

		v := reflect.ValueOf(value)
		if v.Len() != expectedLength {
			return errors.New(fmt.Sprintf("\"%d\" not equal to \"%d\"", v.Len(), expectedLength))
		}
		return nil
	}
}

func Present(expression string) apitest.Assert {
	return func(res *http.Response, req *http.Request) error {
		value, _ := jsonPath(res.Body, expression)
		if isEmpty(value) {
			return errors.New(fmt.Sprintf("value not present for expression: '%s'", expression))
		}
		return nil
	}
}

func NotPresent(expression string) apitest.Assert {
	return func(res *http.Response, req *http.Request) error {
		value, _ := jsonPath(res.Body, expression)
		if !isEmpty(value) {
			return errors.New(fmt.Sprintf("value present for expression: '%s'", expression))
		}
		return nil
	}
}

func Matches(expression string, regexp string) apitest.Assert {
	return func(res *http.Response, req *http.Request) error {
		pattern, err := regex.Compile(regexp)
		if err != nil {
			return errors.New(fmt.Sprintf("invalid pattern: '%s'", regexp))
		}
		value, _ := jsonPath(res.Body, expression)
		if value == nil {
			return errors.New(fmt.Sprintf("no match for pattern: '%s'", expression))
		}
		kind := reflect.ValueOf(value).Kind()
		switch kind {
		case reflect.Bool,
			reflect.Int,
			reflect.Int8,
			reflect.Int16,
			reflect.Int32,
			reflect.Int64,
			reflect.Uint,
			reflect.Uint8,
			reflect.Uint16,
			reflect.Uint32,
			reflect.Uint64,
			reflect.Uintptr,
			reflect.Float32,
			reflect.Float64,
			reflect.String:
			if !pattern.Match([]byte(fmt.Sprintf("%v", value))) {
				return errors.New(fmt.Sprintf("value '%v' does not match pattern '%v'", value, regexp))
			}
			return nil
		default:
			return errors.New(fmt.Sprintf("unable to match using type: %s", kind.String()))
		}
	}
}

func isEmpty(object interface{}) bool {
	if object == nil {
		return true
	}

	objValue := reflect.ValueOf(object)

	switch objValue.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		return objValue.Len() == 0
	case reflect.Ptr:
		if objValue.IsNil() {
			return true
		}
		deref := objValue.Elem().Interface()
		return isEmpty(deref)
	default:
		zero := reflect.Zero(objValue.Type())
		return reflect.DeepEqual(object, zero.Interface())
	}
}

func jsonPath(reader io.Reader, expression string) (interface{}, error) {
	v := interface{}(nil)
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &v)
	if err != nil {
		return nil, err
	}

	value, err := jsonpath.Get(expression, v)
	if err != nil {
		return nil, err
	}
	return value, nil
}

// courtesy of github.com/stretchr/testify
func includesElement(list interface{}, element interface{}) (ok, found bool) {
	listValue := reflect.ValueOf(list)
	elementValue := reflect.ValueOf(element)
	defer func() {
		if e := recover(); e != nil {
			ok = false
			found = false
		}
	}()

	if reflect.TypeOf(list).Kind() == reflect.String {
		return true, strings.Contains(listValue.String(), elementValue.String())
	}

	if reflect.TypeOf(list).Kind() == reflect.Map {
		mapKeys := listValue.MapKeys()
		for i := 0; i < len(mapKeys); i++ {
			if objectsAreEqual(mapKeys[i].Interface(), element) {
				return true, true
			}
		}
		return true, false
	}

	for i := 0; i < listValue.Len(); i++ {
		if objectsAreEqual(listValue.Index(i).Interface(), element) {
			return true, true
		}
	}
	return true, false
}

func objectsAreEqual(expected, actual interface{}) bool {
	if expected == nil || actual == nil {
		return expected == actual
	}

	exp, ok := expected.([]byte)
	if !ok {
		return reflect.DeepEqual(expected, actual)
	}

	act, ok := actual.([]byte)
	if !ok {
		return false
	}
	if exp == nil || act == nil {
		return exp == nil && act == nil
	}
	return bytes.Equal(exp, act)
}
