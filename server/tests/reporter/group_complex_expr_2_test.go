package reporter

import (
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/report"
)

func Test_group_complex_expr_2(t *testing.T) {
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
		"202, 0, 9, 467",
		"202, 1, 3, 183")
}
