package expr

import (
	"reflect"

	"github.com/PaesslerAG/gval"
	"github.com/davecgh/go-spew/spew"
)

func GenericFunctions() []gval.Language {
	return []gval.Language{
		gval.Function("coalesce", coalesce),
		gval.Function("isEmpty", isEmpty),
		gval.Function("isNil", isNil),
	}
}

func coalesce(aa ...interface{}) interface{} {
	for _, a := range aa {
		if !isNil(a) {
			return a
		}
	}

	return nil
}

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Slice, reflect.Array, reflect.Ptr, reflect.Map:
		return reflect.ValueOf(i).IsNil()
	}

	return false
}

// empty checks values and slices
func isEmpty(i interface{}) bool {
	if isNil(i) {
		return true
	}

	spew.Dump("A", reflect.ValueOf(i).IsZero())

	switch reflect.TypeOf(i).Kind() {
	case reflect.Slice, reflect.Array, reflect.Ptr, reflect.Map:
		return reflect.ValueOf(i).Len() == 0
	}

	return reflect.ValueOf(i).IsZero()
}

func isInt(v interface{}) bool {
	switch getKind(v) {
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
		return true
	}

	return false
}

func isString(v interface{}) bool {
	return getKind(v) == reflect.String
}

func isFloat(v interface{}) bool {
	switch getKind(v) {
	case reflect.Float32, reflect.Float64:
		return true
	}

	return false
}

func isBool(v interface{}) bool {
	return getKind(v) == reflect.Bool
}

func getKind(v interface{}) reflect.Kind {
	kind := reflect.TypeOf(v).Kind()

	if isSlice(v) {
		kind = reflect.TypeOf(v).Elem().Kind()

	}

	return kind
}

func isSlice(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Slice
}
