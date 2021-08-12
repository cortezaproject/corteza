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

	f = ff[3]
	h.a.Equal(4, f.Size())
	h.a.Equal("weird_filter", f.Name)
	h.a.Equal("first_name<String>, last_name<String>", f.Columns.String())
	checkRows(h, f,
		"Maria, Königsmann",
		"Engel, Loritz",
		"Maria, Krüger",
		"Engel, Kiefer")

	f = ff[4]
	h.a.Equal(12, f.Size())
	h.a.Equal("empty_filter_object", f.Name)
	h.a.Equal("first_name<String>, last_name<String>", f.Columns.String())
}
