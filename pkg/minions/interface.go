package minions

import "reflect"

// IsNil checks if the given interface is truly nil
//
// Due to how interfaces are handled under-the-hood, a simple i == nil may
// not always be ok.
func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}
