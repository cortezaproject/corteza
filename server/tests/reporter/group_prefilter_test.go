package reporter

import (
	"testing"
)

func Test_group_prefilter(t *testing.T) {
	var (
		ctx, h, s = setup(t)
		m, _, dd  = loadScenario(ctx, s, t, h)
		ff        = loadNoErr(ctx, h, m, dd...)
	)

	h.a.Len(ff, 1)
	f := ff[0]
	h.a.Equal(1, f.Size())

	h.a.Equal("by_name<String>, count<Number>, total<Number>", f.Columns.String())

	checkRows(h, f,
		"Maria, 3, 183")
}
