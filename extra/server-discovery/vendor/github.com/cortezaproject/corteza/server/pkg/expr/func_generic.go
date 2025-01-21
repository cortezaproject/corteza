package expr

import (
	"fmt"
	"github.com/PaesslerAG/gval"
	"github.com/cortezaproject/corteza/server/pkg/gvalfnc"
	"reflect"
)

type (
	empty interface {
		IsEmpty() bool
	}
)

func GenericFunctions() []gval.Language {
	return []gval.Language{
		gval.Function("coalesce", coalesce),
		gval.Function("isEmpty", isEmpty),
		gval.Function("isNil", isNil),
		gval.Function("length", length),
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

func length(i interface{}) int {
	if isEmpty(i) {
		return 0
	}

	switch reflect.TypeOf(i).Kind() {
	case reflect.Slice, reflect.Array, reflect.Ptr, reflect.Map, reflect.String:
		return reflect.ValueOf(i).Len()
	}

	return 0
}

func isNil(i interface{}) bool {
	_, isTyped := i.(TypedValue)

	if isTyped {
		switch i.(type) {
		case *Duration, *DateTime, *Bytes, *Array, *Integer:
			i = i.(TypedValue).Get()
		}
	}

	return gvalfnc.IsNil(i)
}

// empty checks values and slices
func isEmpty(i interface{}) bool {
	if isNil(i) {
		return true
	}

	if c, ok := i.(empty); ok {
		return c.IsEmpty()
	}

	switch reflect.TypeOf(i).Kind() {
	case reflect.Slice, reflect.Array, reflect.Map:
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
	return reflect.TypeOf(v).Kind() == reflect.Slice || reflect.TypeOf(v).Kind() == reflect.Array
}

func isMap(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Map
}

// toArray removes expr types (if wrapped) and checks if the variable is slice
// internal only
func toSlice(vv interface{}) (interface{}, error) {
	vv = UntypedValue(vv)

	if !isSlice(vv) {
		return nil, fmt.Errorf("unexpected type: %T, expecting slice", vv)
	}

	return vv, nil
}
