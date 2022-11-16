package reporter

import (
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/report"
)

func Test_group_complex_expr(t *testing.T) {
	var (
		ctx, h, s = setup(t)
		m, _, dd  = loadScenario(ctx, s, t, h)
		ff        []*report.Frame
		def       = dd[0]
	)

	ff = loadNoErr(ctx, h, m, def)
	h.a.Len(ff, 1)
	r := ff[0]
	h.a.Equal(9, r.Size())

	h.a.Equal("by_year<Number>, by_name<String>, count<Number>, total<Number>", r.Columns.String())

	checkRows(h, ff[0],
		"2020, Engel, 1, 97",
		"2021, Engel, 2, 82",
		"2021, Manu, 1, 61",
		"2020, Maria, 2, 16",
		"2021, Maria, 1, 23",
		"2020, Sascha, 1, 38",
		"2021, Sigi, 1, 67",
		"2020, Ulli, 1, 14",
		"2021, Ulli, 2, 108")
}
