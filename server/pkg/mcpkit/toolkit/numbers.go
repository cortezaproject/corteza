package toolkit

import (
	"fmt"
	"strings"

	"github.com/crusttech/human/server/pkg/j7s"
)

// RejectUnsafeNumbers refuses a decoded argument holding a number too large
// for a JSON number to carry exactly. Human IDs are such numbers, so one sent
// unquoted has already been rounded to a different ID by the time it is read.
func RejectUnsafeNumbers(key string, v any) error {
	path, found := j7s.UnsafeNumber(v)
	if !found {
		return nil
	}

	switch {
	case path == "":
		path = key
	case key != "" && !strings.HasPrefix(path, "["):
		path = key + "." + path
	default:
		path = key + path
	}

	return WithCode(fmt.Errorf(
		"%s is a number too large to be exact in JSON and has already been rounded; "+
			"send IDs and other large numbers as strings", path), CodeInvalidArgument, "")
}
