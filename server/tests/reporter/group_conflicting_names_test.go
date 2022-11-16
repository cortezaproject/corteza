package reporter

import (
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/report"
)

func Test_group_conflicting_names(t *testing.T) {
	var (
		ctx, h, s = setup(t)
		m, _, dd  = loadScenarioOwnDM(ctx, s, t, h)
		ff        []*report.Frame
		def       = dd[0]
	)

	ff = loadNoErr(ctx, h, m, def)
	h.a.Len(ff, 1)
	r := ff[0]
	h.a.Equal(3, r.Size())

	h.a.Equal("by_group<Select>, count<Number>", r.Columns.String())

	checkRows(h, ff[0],
		"g1, 2",
		"g2, 1",
		"g3, 1")
}
