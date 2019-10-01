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

// ParseUInt64s parses a slice of strings into a slice of uint64s
func ParseUInt64s(ss []string) []uint64 {
	uu := make([]uint64, len(ss))
	for i, s := range ss {
		uu[i] = ParseUInt64(s)
	}

	return uu
}
