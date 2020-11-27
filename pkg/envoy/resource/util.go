package resource

import (
	"fmt"
	"strconv"
)

// fn converts identifier values (string, fmt.Stringer, uint64) to string slice
//
// Each value is checked and should not be empty or zero
func identifiers(ii ...interface{}) []string {
	ss := make([]string, 0, len(ii))

	for _, i := range ii {
		switch c := i.(type) {
		case uint64:
			if c == 0 {
				continue
			}

			ss = append(ss, strconv.FormatUint(c, 10))

		case fmt.Stringer:
			if c.String() == "" {
				continue
			}

			ss = append(ss, c.String())

		case string:
			if c == "" {
				continue
			}

			ss = append(ss, c)
		}
	}

	return ss
}
