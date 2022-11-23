package handle

import (
	"regexp"
	"strings"
)

var (
	validHandle  = regexp.MustCompile(`^[A-Za-z][0-9A-Za-z_\-.]*[A-Za-z0-9]$`)
	invalidChars = regexp.MustCompile(`[^0-9A-Za-z_\-.]+`)
)

func IsValid(s string) bool {
	return s == "" || (len(s) >= 2 && validHandle.MatchString(s))
}

// Cast transforms candidates to find a valid (non-empty) handle
func Cast(check func(string) bool, candidates ...string) (handle string, ok bool) {
	ok = true

	for _, c := range candidates {
		if c == "" {
			continue
		}

		// Capitalize
		handle = strings.ReplaceAll(c[:1]+strings.Title(c)[1:], " ", "")
		handle = invalidChars.ReplaceAllString(handle, "")

		if handle == "" {
			continue
		}

		if IsValid(handle) && (check == nil || check(handle)) {
			return
		}
	}

	return "", false
}
