package reporter

import (
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/report"
)

func Test_load_filtering(t *testing.T) {
	var (
		ctx, h, s = setup(t)
		m, _, dd  = loadScenario(ctx, s, t, h)
		ff        = loadNoErr(ctx, h, m, dd...)
	)

	h.a.Len(ff, 1)
	f := ff[0]
	// 3xMaria + 3xUlli + 1xSpecht
	h.a.Equal(7, f.Size())

	h.a.Equal("id<Record>, first_name<String>, last_name<String>, number_of_numbers<Number>", f.Columns.String())
	f.WalkRows(func(i int, r report.FrameRow) error {
		for _, c := range r {
			h.a.NotNil(c)
		}

		return nil
	})
}
