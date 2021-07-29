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
}

func Test2004_grouping_paging(t *testing.T) {
	var (
		ctx, h, s = setup(t)
		m, _, dd  = loadScenario(ctx, s, t, h)
		ff        []*report.Frame
		f         *report.Frame
		def       = dd[0]
	)

	// ^ going up ^
	ff = loadNoErr(ctx, h, m, def)
	h.a.Len(ff, 1)
	f = ff[0]
	h.a.NotNil(f.Paging)
	h.a.NotNil(f.Paging.NextPage)
	h.a.Nil(f.Paging.PrevPage)
	h.a.Equal(4, f.Size())
	checkRows(h, f,
		"Engel, 3, 179",
		"Manu, 1, 61",
		"Maria, 3, 183",
		"Sascha, 1, 38")

	def.Paging.PageCursor = f.Paging.NextPage
	ff = loadNoErr(ctx, h, m, def)
	h.a.Len(ff, 1)
	f = ff[0]
	h.a.NotNil(f.Paging)
	h.a.Nil(f.Paging.NextPage)
	h.a.NotNil(f.Paging.PrevPage)
	h.a.Equal(2, f.Size())
	checkRows(h, f,
		"Sigi, 1, 67",
		"Ulli, 3, 122")

	// v going down v
	def.Paging.PageCursor = f.Paging.PrevPage
	ff = loadNoErr(ctx, h, m, def)
	h.a.Len(ff, 1)
	f = ff[0]
	h.a.NotNil(f.Paging)
	h.a.NotNil(f.Paging.NextPage)
	h.a.Nil(f.Paging.PrevPage)
	h.a.Equal(4, f.Size())
	checkRows(h, f,
		"Engel, 3, 179",
		"Manu, 1, 61",
		"Maria, 3, 183",
		"Sascha, 1, 38")

}
