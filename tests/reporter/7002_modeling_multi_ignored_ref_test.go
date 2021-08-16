package reporter

import (
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/report"
)

func Test7002_modeling_multi_ignored_ref(t *testing.T) {
	var (
		ctx, h, s = setup(t)
		m, _, dd  = loadScenario(ctx, s, t, h)
		ff        = loadNoErrMulti(ctx, h, m, dd...)
		f         *report.Frame
	)

	h.a.Len(ff, 2)

	f = ff[0]
	h.a.Equal("f1", f.Name)
	h.a.Equal("first_name<String>, last_name<String>", f.Columns.String())
	h.a.Len(f.Rows, 3)

	f = ff[1]
	h.a.Equal("f2", f.Name)
	h.a.Equal("first_name<String>, last_name<String>", f.Columns.String())
	h.a.Len(f.Rows, 1)
}
