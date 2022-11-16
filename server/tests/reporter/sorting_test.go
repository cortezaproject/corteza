package reporter

import (
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/report"
)

func Test_sort_empty(t *testing.T) {
	var (
		ctx, h, s = setup(t)
		m, _, dd  = loadScenario(ctx, s, t, h)
		ff        = loadNoErrMulti(ctx, h, m, dd...)
		f         *report.Frame
	)

	h.a.Len(ff, 1)

	f = ff[0]
	h.a.Len(f.Rows, 12)
}
