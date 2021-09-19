package reporter

import (
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/report"
)

func Test_load_basic(t *testing.T) {
	var (
		ctx, h, s = setup(t)
		m, _, dd  = loadScenario(ctx, s, t, h)
		ff        = loadNoErr(ctx, h, m, dd...)
	)

	// size
	h.a.Len(ff, 1)
	f := ff[0]
	h.a.Equal(12, f.Size())

	h.a.Equal("id<Record>, first_name<String>, last_name<String>, number_of_numbers<Number>", f.Columns.String())
	f.WalkRows(func(i int, r report.FrameRow) error {
		for _, c := range r {
			h.a.NotNil(c)
		}

		return nil
	})
}
