package reporter

import (
	"testing"
)

func Test_filter_eval_order(t *testing.T) {
	var (
		ctx, h, s = setup(t)
		m, _, dd  = loadScenarioOwnDM(ctx, s, t, h)
		ff        = loadNoErrMulti(ctx, h, m, dd...)
	)

	f := ff[0]
	h.a.Equal(4, f.Size())
	h.a.Equal("c1", f.Name)
	h.a.Equal("first_name<String>, numbers<Number>", f.Columns.String())
	checkRows(h, f,
		"test1, 11",
		"test1, 11",
		"test2, 11",
		"test1, 14")

	f = ff[1]
	h.a.Equal(4, f.Size())
	h.a.Equal("c2", f.Name)
	h.a.Equal("first_name<String>, numbers<Number>", f.Columns.String())
	checkRows(h, f,
		"test1, 11",
		"test1, 11",
		"test2, 11",
		"test1, 14")

}
