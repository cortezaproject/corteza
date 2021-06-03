package slice

import (
	"reflect"
	"sort"

	"github.com/spf13/cast"
)

// Keys returns sorted map keys
//
// If input is not a map it will return an empty slice
func Keys(m interface{}) (kk []string) {
	v := reflect.ValueOf(m)
	if v.Kind() == reflect.Map {
		kk = make([]string, 0, v.Len())
		for _, kval := range v.MapKeys() {
			if k := cast.ToString(kval.Interface()); k != "" {
				kk = append(kk, k)
			}
		}
	}

	sort.Strings(kk)
	return
}
