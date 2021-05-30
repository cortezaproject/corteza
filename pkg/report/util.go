package report

import (
	"reflect"
	"strings"
)

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}

func dimensionOf(k string) string {
	pp := strings.Split(k, ".")
	if len(pp) < 2 {
		return ""
	}

	return pp[0]
}

func columnOf(k string) string {
	if k == columnWildcard {
		return k
	}

	pp := strings.Split(k, ".")
	if len(pp) < 2 {
		return ""
	}

	return strings.Join(pp[1:], ".")
}
