package slice

func HasUint64(ss []uint64, s uint64) bool {
	for i := range ss {
		if ss[i] == s {
			return true
		}
	}

	return false
}
