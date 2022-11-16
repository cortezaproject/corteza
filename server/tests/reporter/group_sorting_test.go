package reporter

import (
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/report"
)

func Test_group_sorting(t *testing.T) {
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
