package reporter

import (
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/report"
)

func Test_filter_ast(t *testing.T) {
	var (
		ctx, h, s = setup(t)
		m, _, dd  = loadScenario(ctx, s, t, h)
		ff        = loadNoErrMulti(ctx, h, m, dd...)
		f         *report.Frame
	)

	h.a.Len(ff, 1)

	f = ff[0]
	h.a.Len(f.Rows, 5)
	h.a.Equal("first_name<String>, last_name<String>", f.Columns.String())
	checkRows(h, f,
		"Maria, Königsmann",
		"Maria, Spannagel",
		"Engel, Kempf",
		"Maria, Krüger",
		"Engel, Kiefer")
}
