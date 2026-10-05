package j7s

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// MaxExactNumber is 2^53, the magnitude from which a JSON number decoded as
// float64 (or read by JavaScript) no longer holds every integer exactly.
const MaxExactNumber = 1 << 53

// UnsafeNumber returns the path of the first float64 in v whose magnitude is
// at least MaxExactNumber — a number that was rounded when it was decoded.
func UnsafeNumber(v any) (path string, found bool) {
	return unsafeNumber(v, "")
}

func unsafeNumber(v any, path string) (string, bool) {
	switch c := v.(type) {
	case float64:
		if math.Abs(c) >= MaxExactNumber {
			return path, true
		}
	case map[string]any:
		for k, item := range c {
			p := k
			if path != "" {
				p = path + "." + k
			}
			if found, ok := unsafeNumber(item, p); ok {
				return found, true
			}
		}
	case []any:
		for i, item := range c {
			if found, ok := unsafeNumber(item, fmt.Sprintf("%s[%d]", path, i)); ok {
				return found, true
			}
		}
	}

	return "", false
}

// ExactNumbers replaces each json.Number in v (as decoded with UseNumber): an
// integer beyond MaxExactNumber becomes the string of its digits, anything
// else a float64.
func ExactNumbers(v any) any {
	switch c := v.(type) {
	case json.Number:
		s := c.String()
		if !strings.ContainsAny(s, ".eE") {
			if n, err := strconv.ParseInt(s, 10, 64); err == nil && n > -MaxExactNumber && n < MaxExactNumber {
				return float64(n)
			}
			if _, err := strconv.ParseUint(strings.TrimPrefix(s, "-"), 10, 64); err == nil {
				return s
			}
		}
		f, _ := c.Float64()
		return f
	case map[string]any:
		for k, item := range c {
			c[k] = ExactNumbers(item)
		}
	case []any:
		for i, item := range c {
			c[i] = ExactNumbers(item)
		}
	}

	return v
}
