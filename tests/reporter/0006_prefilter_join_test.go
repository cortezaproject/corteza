package reporter

import (
	"testing"
)

func Test0006_prefilter_join(t *testing.T) {
	var (
		ctx, h, s = setup(t)
		m, _, dd  = loadScenario(ctx, s, t, h)
		ff        = loadNoErr(ctx, h, m, dd...)
	)

	h.a.Len(ff, 2)
	f := ff[0]

	h.a.Equal(1, f.Size())
	h.a.Equal("id<Record>, join_key<String>, first_name<String>, last_name<String>", f.Columns.String())
	checkRows(h, f,
		", Maria_Königsmann, Maria, Königsmann")

	f = ff[1]

	h.a.Equal(2, f.Size())
	h.a.Equal("id<Record>, usr<String>, name<String>, type<Select>, cost<Number>, time_spent<Number>", f.Columns.String())
	checkRows(h, f,
		", Maria_Königsmann, u1 j1, a, 10, 2",
		", Maria_Königsmann, u1 j5, a, 4, 4")
}
