package internal

import "strings"

func CamelCase(pp ...string) (out string) {
	for i, p := range pp {
		if i > 0 && len(p) > 1 {
			p = strings.ToUpper(p[:1]) + p[1:]
		}

		out = out + p
	}

	return out
}
