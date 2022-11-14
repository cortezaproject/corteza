package slice

func IntersectStrings(a []string, b []string) []string {
	var (
		out = make([]string, 0, len(a)+len(b))
		ah  = ToStringBoolMap(a)
	)

	for i := 0; i < len(b); i++ {
		if ah[b[i]] {
			out = append(out, b[i])
		}
	}

	return out
}

func ToStringBoolMap(s []string) (h map[string]bool) {
	h = make(map[string]bool)
	for i := 0; i < len(s); i++ {
		h[s[i]] = true
	}

	return
}

func ToUint64BoolMap(u []uint64) (h map[uint64]bool) {
	h = make(map[uint64]bool)
	for i := 0; i < len(u); i++ {
		h[u[i]] = true
	}

	return
}

func HasString(ss []string, s string) bool {
	for i := range ss {
		if ss[i] == s {
			return true
		}
	}

	return false
}

// RemoveString removes one or more strings form the input slice
func PluckString(ss []string, ff ...string) (o []string) {
	if len(ff) == 0 {
		return ss
	}

	f := ToStringBoolMap(ff)
	o = make([]string, 0, len(ss))

	for _, s := range ss {
		if !f[s] {
			// remove from the list
			o = append(o, s)
		}
	}

	return o
}
