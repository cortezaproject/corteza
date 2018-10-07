package payload

import (
	"strconv"
)

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

// ParseUInt64 parses an string to uint64
func ParseUInt64(s string) uint64 {
	if s == "" {
		return 0
	}
	i, _ := strconv.ParseUint(s, 10, 64)
	return i
}
