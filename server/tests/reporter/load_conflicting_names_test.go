package reporter

import (
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/report"
)

func Test_load_conflicting_names(t *testing.T) {
	var (
		ctx, h, s = setup(t)
		m, _, dd  = loadScenarioOwnDM(ctx, s, t, h)
		ff        []*report.Frame
		def       = dd[0]
	)

	ff = loadNoErr(ctx, h, m, def)
	h.a.Len(ff, 1)
	r := ff[0]
	h.a.Equal(4, r.Size())

	h.a.Equal("name<String>, group<Select>", r.Columns.String())

	checkRows(h, ff[0],
		"item_1_g1, g1",
		"item_2_g1, g1",
		"item_3_g2, g2",
		"item_4_g3, g3")
}
