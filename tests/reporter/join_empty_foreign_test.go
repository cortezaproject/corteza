package reporter

import (
	"testing"
)

func Test_join_empty_foreign(t *testing.T) {
	var (
		ctx, h, s = setup(t)
		m, _, dd  = loadScenario(ctx, s, t, h)
		ff        = loadNoErr(ctx, h, m, dd...)
	)

	h.a.Len(ff, 1)
	f := ff[0]

	h.a.Equal(0, f.Size())
}
