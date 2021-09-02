package reporter

import (
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/report"
)

func Test2001_grouping_basic(t *testing.T) {
	var (
		ctx, h, s = setup(t)
		m, _, dd  = loadScenario(ctx, s, t, h)
		ff        []*report.Frame
		def       = dd[0]
	)

	ff = loadNoErr(ctx, h, m, def)
	h.a.Len(ff, 1)
	r := ff[0]
	h.a.Equal(6, r.Size())

	h.a.Equal("by_name<String>, count<Number>, total<Number>", r.Columns.String())

	checkRows(h, ff[0],
		"Engel, 3, 179",
		"Manu, 1, 61",
		"Maria, 3, 183",
		"Sascha, 1, 38",
		"Sigi, 1, 67",
		"Ulli, 3, 122")
}

func Test2002_grouping_filtering(t *testing.T) {
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

	h.a.Equal("by_name<String>, count<Number>, total<Number>", r.Columns.String())

	checkRows(h, ff[0],
		"Manu, 1, 61",
		"Sigi, 1, 67")
}

func Test2003_grouping_sorting(t *testing.T) {
	var (
		ctx, h, s = setup(t)
		m, _, dd  = loadScenario(ctx, s, t, h)
		ff        []*report.Frame
		def       = dd[0]
	)

	ff = loadNoErr(ctx, h, m, def)
	h.a.Len(ff, 1)
	r := ff[0]
	h.a.Equal(6, r.Size())

	h.a.Equal("by_name<String>, count<Number>, total<Number>", r.Columns.String())

	checkRows(h, ff[0],
		"Engel, 3, 179",
		"Maria, 3, 183",
		"Ulli, 3, 122",
		"Manu, 1, 61",
		"Sascha, 1, 38",
		"Sigi, 1, 67")

	h.a.Equal("count DESC, by_name", r.Sort.String())
}

func Test2004_grouping_paging(t *testing.T) {
	t.Skip("@todo how can we support paging for groupped data? We need to assure something unique")
}
