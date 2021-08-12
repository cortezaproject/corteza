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

	t.Run(ff[0].Name, func(t *testing.T) {
		f := ff[0]
		h.a.Equal(12, f.Size())
		h.a.Equal("empty_filter", f.Name)
		h.a.Equal("first_name<String>, last_name<String>", f.Columns.String())
	})

	t.Run(ff[1].Name, func(t *testing.T) {
		f := ff[1]
		h.a.Equal(0, f.Size())
		h.a.Equal("false_filter", f.Name)
		h.a.Equal("first_name<String>, last_name<String>", f.Columns.String())
	})

	t.Run(ff[2].Name, func(t *testing.T) {
		f := ff[2]
		h.a.Equal(12, f.Size())
		h.a.Equal("true_filter", f.Name)
		h.a.Equal("first_name<String>, last_name<String>", f.Columns.String())
	})

	t.Run(ff[3].Name, func(t *testing.T) {
		f := ff[3]
		h.a.Equal(4, f.Size())
		h.a.Equal("weird_filter", f.Name)
		h.a.Equal("first_name<String>, last_name<String>", f.Columns.String())
		checkRows(h, f,
			"Maria, Königsmann",
			"Engel, Loritz",
			"Maria, Krüger",
			"Engel, Kiefer")

	})

	t.Run(ff[4].Name, func(t *testing.T) {
		f := ff[4]
		h.a.Equal(12, f.Size())
		h.a.Equal("empty_filter_object", f.Name)
		h.a.Equal("first_name<String>, last_name<String>", f.Columns.String())
	})

	t.Run(ff[5].Name, func(t *testing.T) {
		f := ff[5]
		h.a.Equal(1, f.Size())
		h.a.Equal("wildcard_number_value", f.Name)
		h.a.Equal("first_name<String>, last_name<String>", f.Columns.String())
		checkRows(h, f,
			"Engel, Loritz")
	})
}
