package deinterfacer

import (
	"strconv"

	"github.com/pkg/errors"
)

// each() tries to resolve interface{} into map[interface{}]interface{} or []interface{}
// and calls fn() for each entry
func Each(i interface{}, fn func(int, string, interface{}) error) (err error) {
	if kv, ok := i.(map[interface{}]interface{}); ok {
		for k, v := range kv {
			if key, ok := k.(string); !ok {
				err = errors.WithStack(errors.Errorf("unsupported key type: %T (%+v)", k, k))
			} else {
				err = fn(-1, key, v)
			}

			if err != nil {
				return
			}
		}
	} else if kv, ok := i.(map[string]interface{}); ok {
		for k, v := range kv {
			err = fn(-1, k, v)

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

func Simplify(i interface{}) interface{} {
	if IsMap(i) {
		kv := map[string]interface{}{}
		Each(i, func(_ int, k string, v interface{}) error {
			kv[k] = Simplify(v)
			return nil
		})

		return kv
	} else if IsSlice(i) {
		s := make([]interface{}, len(i.([]interface{})))

		Each(i, func(n int, _ string, v interface{}) error {
			s[n] = Simplify(v)
			return nil
		})

		return s
	} else {
		return i
	}
}

// KVsetString assigns value (if exists and of type string) from src to destination
func KVsetString(dst *string, key string, src interface{}, def ...string) {
	if kv, ok := src.(map[interface{}]interface{}); !ok {
		return
	} else if val, ok := kv[key]; ok {
		*dst = ToString(val, def...)
	}
}

// ToString assigns value (if exists and of type string) from src to destination
func ToString(val interface{}, def ...string) string {
	if str, ok := val.(string); ok {
		return str
	} else if len(def) > 0 {
		return def[0]
	}
	return ""
}

// ToInt assigns value (if exists and of type int) from src to destination
func ToInt(val interface{}, def ...int) int {
	if i, ok := val.(int); ok {
		return i
	} else if len(def) > 0 {
		return def[0]
	}
	return 0
}

// ToInt assigns value (if exists and of type int) from src to destination
func ToUint64(val interface{}, def ...uint64) uint64 {
	if i, ok := val.(uint64); ok {
		return i
	} else if s, ok := val.(string); ok {
		if i, _ := strconv.ParseUint(s, 10, 64); i > 0 {
			return i
		}
	}

	if len(def) > 0 {
		return def[0]
	}
	return 0
}

// ToBool assigns value (if exists and of type bool) from src to destination
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

func ToSliceOfStringToInterfaceMap(def interface{}) []map[string]interface{} {
	items := make([]map[string]interface{}, 0)
	_ = Each(def, func(_ int, _ string, def interface{}) error {
		item := make(map[string]interface{})
		_ = Each(def, func(_ int, k string, v interface{}) error {
			item[k] = Simplify(v)
			return nil
		})

		items = append(items, item)
		return nil
	})

	return items
}

// ToInt assigns values (if exists and of type ints) from src to destination
func ToInts(i interface{}, def ...int) (out []int) {
	var (
		ii  []interface{}
		num int
		ok  bool
		n   int
	)

	if ii, ok = i.([]interface{}); ok {
		out = make([]int, len(ii))
		for n, i = range ii {
			if num, ok = i.(int); ok {
				out[n] = num
			}
		}
	} else if num, ok = i.(int); ok {
		return []int{num}
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
