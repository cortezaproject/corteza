package cast2

// Anys converts any kinds of values to a []any slice
func Anys[C any](in ...C) (out []any) {
	out = make([]any, len(in))

	for i := range in {
		out[i] = in[i]
	}

	return
}
