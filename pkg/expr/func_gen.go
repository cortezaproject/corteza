package expr

import (
	"github.com/PaesslerAG/gval"
	"reflect"
)

func GenericFunctions() []gval.Language {
	return []gval.Language{
		gval.Function("coalesce", coalesce),
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
