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
