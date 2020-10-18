package store

import (
	"fmt"
	"strings"
)

// PreprocessValue value preprocessor logic to modify input value before using it in condition filters and other logic
//
// We intentionally panic in case of an error here; misuse can only happen with in case of misconfiguration.
func PreprocessValue(val interface{}, p string) interface{} {
	switch p {
	case "":
		return val
	case "lower":
		if str, ok := val.(string); ok {
			return strings.ToLower(str)
		}
		panic(fmt.Sprintf("preprocessor %q not compatible with type %T (value: %v)", p, val, val))

	default:
		panic(fmt.Sprintf("unknown preprocessor %q used for value %v", p, val))
	}
}
