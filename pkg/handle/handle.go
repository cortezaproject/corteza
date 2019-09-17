package handle

import (
	"regexp"
)

var (
	c = regexp.MustCompile(`^[A-Za-z][0-9A-Za-z_\-.]*[A-Za-z0-9]$`)
)

func IsValid(s string) bool {
	return s == "" || (len(s) >= 2 && c.MatchString(s))
}
