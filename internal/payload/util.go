package payload

import (
	"regexp"
	"strconv"
)

var truthy = regexp.MustCompile("^\\s*(t(rue)?|y(es)?|1)\\s*$")

func Uint64toa(i uint64) string {
	return strconv.FormatUint(i, 10)
}

func Uint64stoa(uu []uint64) []string {
	ss := make([]string, len(uu))
	for i, u := range uu {
		ss[i] = Uint64toa(u)
	}

	return ss
}

//// parseInt64 parses an string to int64
//func parseInt64(s string) int64 {
//	if s == "" {
//		return 0
//	}
//	i, _ := strconv.ParseInt(s, 10, 64)
//
//	return i
//}
//
// parseUInt64 parses an string to uint64
func ParseUInt64(s string) uint64 {
	if s == "" {
		return 0
	}
	i, _ := strconv.ParseUint(s, 10, 64)
	return i
}

//// parseUInt64 parses an string to uint64
//func parseBool(s string) bool {
//	return truthy.MatchString(strings.ToLower(s))
//}
//
//// is checks if string s is contained in matches
//func is(s string, matches ...string) bool {
//	for _, v := range matches {
//		if s == v {
//			return true
//		}
//	}
//	return false
//}
