package reporter

import (
	"testing"
)

func Test9007_filtering_specials(t *testing.T) {
	var (
		ctx, h, s = setup(t)
		m, _, dd  = loadScenario(ctx, s, t, h)
		ff        = loadNoErrMulti(ctx, h, m, dd...)
	)

	h.a.Len(ff, 3)
	f := ff[0]
	h.a.Equal(12, f.Size())
	h.a.Equal("empty_filter", f.Name)
	h.a.Equal("first_name<String>, last_name<String>", f.Columns.String())

	f = ff[1]
	h.a.Equal(0, f.Size())
	h.a.Equal("false_filter", f.Name)
	h.a.Equal("first_name<String>, last_name<String>", f.Columns.String())

	f = ff[2]
	h.a.Equal(12, f.Size())
	h.a.Equal("true_filter", f.Name)
	h.a.Equal("first_name<String>, last_name<String>", f.Columns.String())
}
