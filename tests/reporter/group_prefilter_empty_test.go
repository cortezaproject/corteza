package reporter

import (
	"testing"
)

func Test_group_prefilter_empty(t *testing.T) {
	var (
		ctx, h, s = setup(t)
		m, _, dd  = loadScenario(ctx, s, t, h)
		ff        = loadNoErr(ctx, h, m, dd...)
	)

	h.a.Len(ff, 1)
	f := ff[0]
	h.a.Equal(0, f.Size())

	h.a.Equal("by_name<String>, count<Number>, total<Number>", f.Columns.String())
}
