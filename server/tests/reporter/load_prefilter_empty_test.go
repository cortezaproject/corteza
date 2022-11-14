package reporter

import (
	"testing"
)

func Test_load_prefilter_empty(t *testing.T) {
	var (
		ctx, h, s = setup(t)
		m, _, dd  = loadScenario(ctx, s, t, h)
		ff        = loadNoErr(ctx, h, m, dd...)
	)

	h.a.Len(ff, 1)
	f := ff[0]
	h.a.Equal(0, f.Size())

	h.a.Equal("id<Record>, first_name<String>, last_name<String>, number_of_numbers<Number>", f.Columns.String())
}
