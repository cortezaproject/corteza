package jsonpath

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"io"
	"io/ioutil"
	"reflect"
	"strings"
)

func Contains(expression string, expected interface{}, data io.Reader) error {
	value, err := JsonPath(data, expression)
	if err != nil {
		return err
	}
	ok, found := IncludesElement(value, expected)
	if !ok {
		return errors.New(fmt.Sprintf("\"%s\" could not be applied builtin len()", expected))
	}
	if !found {
		return errors.New(fmt.Sprintf("\"%s\" does not contain \"%s\"", value, expected))
	}
	return nil
}

func Equal(expression string, expected interface{}, data io.Reader) error {
	value, err := JsonPath(data, expression)
	if err != nil {
		return err
	}
	if !ObjectsAreEqual(value, expected) {
		return errors.New(fmt.Sprintf("\"%s\" not equal to \"%s\"", value, expected))
	}
	return nil
}

func NotEqual(expression string, expected interface{}, data io.Reader) error {
	value, err := JsonPath(data, expression)
	if err != nil {
		return err
	}

	if ObjectsAreEqual(value, expected) {
		return errors.New(fmt.Sprintf("\"%s\" value is equal to \"%s\"", expression, expected))
	}
	return nil
}

func Length(expression string, expectedLength int, data io.Reader) error {
	value, err := JsonPath(data, expression)
	if err != nil {
		return err
	}

	v := reflect.ValueOf(value)
	if v.Len() != expectedLength {
		return errors.New(fmt.Sprintf("\"%d\" not equal to \"%d\"", v.Len(), expectedLength))
	}
	return nil
}

func GreaterThan(expression string, minimumLength int, data io.Reader) error {
	value, err := JsonPath(data, expression)
	if err != nil {
		return err
	}

	v := reflect.ValueOf(value)
	if v.Len() < minimumLength {
		return errors.New(fmt.Sprintf("\"%d\" is greater than \"%d\"", v.Len(), minimumLength))
	}
	return nil
}

func LessThan(expression string, maximumLength int, data io.Reader) error {
	value, err := JsonPath(data, expression)
	if err != nil {
		return err
	}

	v := reflect.ValueOf(value)
	if v.Len() > maximumLength {
		return errors.New(fmt.Sprintf("\"%d\" is less than \"%d\"", v.Len(), maximumLength))
	}
	return nil
}

func Present(expression string, data io.Reader) error {
	value, _ := JsonPath(data, expression)
	if isEmpty(value) {
		return errors.New(fmt.Sprintf("value not present for expression: '%s'", expression))
	}
	return nil
}

func NotPresent(expression string, data io.Reader) error {
	value, _ := JsonPath(data, expression)
	if !isEmpty(value) {
		return errors.New(fmt.Sprintf("value present for expression: '%s'", expression))
	}
	return nil
}

func JsonPath(reader io.Reader, expression string) (interface{}, error) {
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
		return nil, fmt.Errorf("evaluating '%s' resulted in error: '%s'", expression, err)
	}
	return value, nil
}

// courtesy of github.com/stretchr/testify
func IncludesElement(list interface{}, element interface{}) (ok, found bool) {
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
			if ObjectsAreEqual(mapKeys[i].Interface(), element) {
				return true, true
			}
		}
		return true, false
	}

	for i := 0; i < listValue.Len(); i++ {
		if ObjectsAreEqual(listValue.Index(i).Interface(), element) {
			return true, true
		}
	}
	return true, false
}

func ObjectsAreEqual(expected, actual interface{}) bool {
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