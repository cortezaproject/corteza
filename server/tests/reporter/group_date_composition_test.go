package reporter

import (
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/report"
)

func Test_group_date_composition(t *testing.T) {
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

	h.a.Equal("by_dob<String>, count<Number>", r.Columns.String())

	checkRows(h, ff[0],
		"2020-01, 5",
		"2021-01, 7")
}
