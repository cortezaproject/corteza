package reporter

import (
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/report"
)

func Test7001_modeling_multi(t *testing.T) {
	var (
		ctx, h, s = setup(t)
		m, _, dd  = loadScenario(ctx, s, t, h)
		ff        = loadNoErrMulti(ctx, h, m, dd...)
		f         *report.Frame
	)

	h.a.Len(ff, 2)

	f = ff[0]
	h.a.Equal("users", f.Source)
	h.a.Equal("first_name<String>, last_name<String>", f.Columns.String())
	checkRows(h, f,
		"Engel, Loritz",
		"Engel, Kempf",
		"Engel, Kiefer",
		"Manu, Specht",
		"Maria, Königsmann",
		"Maria, Spannagel",
		"Maria, Krüger",
		"Sascha, Jans",
		"Sigi, Goldschmidt",
		"Ulli, Haupt",
		"Ulli, Böhler",
		"Ulli, Förstner")

	f = ff[1]
	h.a.Equal("grouped", f.Source)
	h.a.Equal("by_name<String>, count<Number>, total<Number>", f.Columns.String())
	checkRows(h, f,
		"Engel, 3, 179",
		"Manu, 1, 61",
		"Maria, 3, 183",
		"Sascha, 1, 38",
		"Sigi, 1, 67",
		"Ulli, 3, 122")
}
