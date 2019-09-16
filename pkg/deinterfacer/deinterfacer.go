package deinterfacer

import (
	"github.com/pkg/errors"
)

// each() tries to resolve interface{} into map[interface{}]interface{} or []interface{}
// and calls fn() for each entry
func Each(i interface{}, fn func(int, string, interface{}) error) (err error) {
	if kv, ok := i.(map[interface{}]interface{}); ok {
		for k, v := range kv {
			if key, ok := k.(string); !ok {
				err = errors.WithStack(errors.New("unsupported key type"))
			} else {
				err = fn(-1, key, v)
			}

			if err != nil {
				return
			}
		}
	} else if slice, ok := i.([]interface{}); ok {
		_ = slice
		for index, i := range slice {
			if err = fn(index, "", i); err != nil {
				return
			}
		}
	}

	return
}

// KVsetString assigns value (if exists and of type string) from src to destination
func KVsetString(dst *string, key string, src interface{}, def ...string) {
	if kv, ok := src.(map[interface{}]interface{}); !ok {
		return
	} else if val, ok := kv[key]; ok {
		*dst = ToString(val, def...)
	}
}

// KVsetString assigns value (if exists and of type string) from src to destination
func ToString(val interface{}, def ...string) string {
	if str, ok := val.(string); ok {
		return str
	} else if len(def) > 0 {
		return def[0]
	}
	return ""
}

// KVsetString assigns value (if exists and of type string) from src to destination
func ToBool(val interface{}, def ...bool) bool {
	if b, ok := val.(bool); ok {
		return b
	} else if len(def) > 0 {
		return def[0]
	}
	return false
}

func ToStrings(i interface{}) (out []string) {
	var (
		ii  []interface{}
		str string
		ok  bool
		n   int
	)

	if ii, ok = i.([]interface{}); ok {
		out = make([]string, len(ii))
		for n, i = range ii {
			if str, ok = i.(string); ok {
				out[n] = str
			}
		}
	} else if str, ok = i.(string); ok {
		return []string{str}
	}

	return
}

func IsMap(i interface{}) (ok bool) {
	_, ok = i.(map[interface{}]interface{})
	return
}

func IsSlice(i interface{}) (ok bool) {
	_, ok = i.([]interface{})
	return
}

func IsIterable(i interface{}) (ok bool) {
	return IsMap(i) || IsSlice(i)
}
