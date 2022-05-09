package types

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type (
	// KVDecoder interface for custom decoding logic
	KVDecoder interface {
		DecodeKV(SettingsKV, string) error
	}
)

// DecodeKV converts key-value (SettingsKV) into structs using tags & field names
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
func DecodeKV(kv SettingsKV, dst interface{}, pp ...string) (err error) {
	valueOf := reflect.ValueOf(dst)
	if valueOf.Kind() != reflect.Ptr {
		return errors.New("expecting a pointer, not a value")
	}

	if valueOf.IsNil() {
		return errors.New("nil pointer passed")
	}

	var prefix string
	if len(pp) > 0 {
		// If called with prefix, join string slice + 1 empty string (to ensure tailing dot)
		prefix = strings.Join(append(pp, ""), ".")
	}

	valueOf = valueOf.Elem()

	length := valueOf.NumField()

	for i := 0; i < length; i++ {
		var (
			structField = valueOf.Field(i)
			tags        = make(map[string]bool)
		)

		if !structField.CanSet() {
			continue
		}

		var (
			structFType = valueOf.Type().Field(i)

			// whe nwe use name of the struct field directly, remove upper-case from first letter
			key = prefix + strings.ToLower(structFType.Name[:1]) + structFType.Name[1:]

			tag = structFType.Tag.Get("kv")

			tagFlags []string
		)

		if tag == "-" {
			// Skip fields with kv tags equal to "-"
			continue
		}

		if tag != "" {
			tagFlags = strings.Split(tag, ",")

			if tagFlags[0] != "" {
				// key explicitly set via tag, use that!
				key = prefix + tagFlags[0]
			}

			for f := 1; f < len(tagFlags); f++ {
				// @todo resolve i18n and other flags...
				tags[tagFlags[f]] = true
			}
		}

		var structValue interface{}

		if structField.Kind() == reflect.Ptr {
			structValue = structField.Interface()
		} else {
			structValue = structField.Addr().Interface()
		}

		// Handle custom KVDecoder
		if decodeMethod := reflect.ValueOf(structValue).MethodByName("DecodeKV"); decodeMethod.IsValid() {
			if decode, ok := decodeMethod.Interface().(func(SettingsKV, string) error); !ok {
				panic("invalid DecodeKV() function signature")
			} else if err = decode(kv, key); err != nil {
				return fmt.Errorf("cannot decode settings for %q: %w", key, err)
			} else {
				continue
			}
		}

		if !tags["final"] {
			// Handles structs
			//
			// It calls DecodeKV recursively
			if structField.Kind() == reflect.Struct {
				if err = DecodeKV(kv.Filter(key), structValue, key); err != nil {
					return err
				}

				continue
			}

			// Handles map values
			if structField.Kind() == reflect.Map {
				if structField.IsNil() {
					// allocate new map
					structField.Set(reflect.MakeMap(structField.Type()))
				}

				// cut SettingsKV key prefix and use the rest for the map key
				for k, val := range kv.CutPrefix(key + ".") {
					mapValue := reflect.New(structField.Type().Elem())
					err = val.Unmarshal(mapValue.Interface())
					if err != nil {
						return fmt.Errorf("cannot decode JSON into map for key %q: %w", key, err)
					}

					structField.SetMapIndex(reflect.ValueOf(k), mapValue.Elem())
				}

				continue
			}
		}

		// Native type
		if val, ok := kv[key]; ok {
			// Always use pointer to value
			if val == nil {
				switch structFType.Type.Kind() {
				case reflect.String:
					structField.SetString("")
					continue
				}
			}

			if val.Unmarshal(structField.Addr().Interface()) != nil {
				// Try to get numbers encoded as strings...
				var tmp interface{}
				if err = val.Unmarshal(&tmp); err != nil {
					return fmt.Errorf("could not decode JSON for key %q: %w", key, err)
				}

				var cnv, is = tmp.(string)
				if !is {
					// give up
					continue
				}

				switch structFType.Type.Kind() {
				case reflect.Int, reflect.Int32, reflect.Int64:
					if num, err := strconv.ParseInt(cnv, 10, 64); err == nil {
						structField.SetInt(num)
					}
				case reflect.Uint, reflect.Uint32, reflect.Uint64:
					if num, err := strconv.ParseUint(cnv, 10, 64); err == nil {
						structField.SetUint(num)
					}
				case reflect.Float32, reflect.Float64:
					if num, err := strconv.ParseFloat(cnv, 64); err == nil {
						structField.SetFloat(num)
					}

				}
			}
		}
	}

	return
}
