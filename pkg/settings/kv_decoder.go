package settings

import (
	"errors"
	"reflect"
	"strings"
)

type (
// @todo support Decoder interface
// Decoder interface {
// 	Decode(kv KV, prefix string) error
// }
)

var (
// @todo support Decoder interface
// decoderTyEl = reflect.TypeOf((*Decoder)(nil)).Elem()
)

// Decode converts key-value (KV) into structs using tags & field names
//
// Supports decoding into all scalar types, can handle nested structures and simple maps (1 dim, string as key)
//
// Example:
// key-value pairs:
//   string1: "string"
//   number:  42
//
// struct{
//   String1 string `kv:"string1`
//   Number  int
// }
//
func Decode(kv KV, dst interface{}, pp ...string) (err error) {
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Ptr {
		return errors.New("expecting a pointer, not a value")
	}

	if v.IsNil() {
		return errors.New("nil pointer passed")
	}

	v = v.Elem()

	var prefix string
	if len(pp) > 0 {
		// If called with prefix, join string slice + 1 empty string (to ensure tailing dot)
		prefix = strings.Join(append(pp, ""), ".")
	}

	length := v.NumField()

	for i := 0; i < length; i++ {
		var (
			f = v.Field(i)
			t = v.Type().Field(i)

			key = prefix + strings.ToLower(t.Name[:1]) + t.Name[1:]
			tag = t.Tag.Get("kv")

			tagFlags []string
		)

		if tag == "-" {
			// Skip fields with kv tags equal to "-"
			continue
		}

		// if !f.CanSet() {
		// 	return errors.New("unexpected pointer for field " + t.Name)
		// }

		if tag != "" {
			tagFlags = strings.Split(tag, ",")

			if tagFlags[0] != "" {
				// key explicitly set via tag, use that!
				key = prefix + tagFlags[0]
			}

			for f := 1; f < len(tagFlags); f++ {
				// @todo resolve i18n and other flags...
			}
		}

		// @todo handle Decoder interface
		// if f.Type().Implements(decoderTyEl) {
		// 	result := reflect.ValueOf(&t).MethodByName("Decode").Call([]reflect.Value{
		// 		reflect.ValueOf(kv.Filter(key)),
		// 		reflect.ValueOf(prefix),
		// 	})
		//
		// 	if len(result) != 1 {
		// 		return errors.New("internal error, Decoder signature does not match")
		// 	}
		// }

		// Handles structs
		//
		// It calls Decode recursively
		if f.Kind() == reflect.Struct {
			if err = Decode(kv.Filter(key), f.Addr().Interface(), key); err != nil {
				return
			}

			continue
		}

		// Handles map values
		if f.Kind() == reflect.Map {
			if f.IsNil() {
				// allocate new map
				f.Set(reflect.MakeMap(f.Type()))
			}

			for k, val := range kv.CutPrefix(key + ".") {
				mapValue := reflect.New(f.Type().Elem())
				val.Unmarshal(mapValue.Interface())
				f.SetMapIndex(reflect.ValueOf(k), mapValue.Elem())
			}

			continue
		}

		if val, ok := kv[key]; ok {
			if err = val.Unmarshal(f.Addr().Interface()); err != nil {
				return
			}

		}
	}

	return
}
