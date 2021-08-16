package reporter

import (
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/report"
)

func Test2005_grouping_complex_ast(t *testing.T) {
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
		"2020, Engel, 1, 9.7",
		"2021, Engel, 2, 8.2",
		"2021, Manu, 1, 6.1",
		"2020, Maria, 2, 16",
		"2021, Maria, 1, 2.3",
		"2020, Sascha, 1, 3.8",
		"2021, Sigi, 1, 6.7",
		"2020, Ulli, 1, 1.4",
		"2021, Ulli, 2, 10.8")
}

func Test2006_grouping_complex_expr(t *testing.T) {
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
		"2020, Engel, 1, 9.7",
		"2021, Engel, 2, 8.2",
		"2021, Manu, 1, 6.1",
		"2020, Maria, 2, 16",
		"2021, Maria, 1, 2.3",
		"2020, Sascha, 1, 3.8",
		"2021, Sigi, 1, 6.7",
		"2020, Ulli, 1, 1.4",
		"2021, Ulli, 2, 10.8")
}

func Test2007_grouping_complex_expr_2(t *testing.T) {
	var (
		ctx, h, s = setup(t)
		m, _, dd  = loadScenario(ctx, s, t, h)
		ff        []*report.Frame
		def       = dd[0]
	)

	ff = loadNoErr(ctx, h, m, def)
	h.a.Len(ff, 1)
	r := ff[0]
	h.a.Equal(2, r.Size())

	h.a.Equal("by_year<Number>, is_maria<Boolean>, count<Number>, total<Number>", r.Columns.String())

	checkRows(h, ff[0],
		"202, 0, 9, 46.7",
		"202, 1, 3, 18.3")
}

func Test2008_grouping_complex_ast_expr(t *testing.T) {
	var (
		ctx, h, s = setup(t)
		m, _, dd  = loadScenario(ctx, s, t, h)
		ff        []*report.Frame
		def       = dd[0]
	)

	ff = loadNoErr(ctx, h, m, def)
	h.a.Len(ff, 1)
	r := ff[0]
	h.a.Equal(2, r.Size())

	h.a.Equal("by_year<Number>, is_maria<Boolean>, count<Number>, total<Number>", r.Columns.String())

	checkRows(h, ff[0],
		"202, 0, 9, 46.7",
		"202, 1, 3, 18.3")
}
