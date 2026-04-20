package id

import "strconv"

func Strings(ii ...uint64) []string {
	ss := make([]string, len(ii))
	for i, v := range ii {
		ss[i] = String(v)
	}
	return ss
}

func String(i uint64) string {
	return strconv.FormatUint(i, 10)
}

func Uints(ss ...string) []uint64 {
	uu := make([]uint64, len(ss))
	for i, s := range ss {
		uu[i] = Uint(s)
	}
	return uu
}

func Uint(s string) uint64 {
	if s == "" {
		return 0
	}
	i, _ := strconv.ParseUint(s, 10, 64)
	return i
}
