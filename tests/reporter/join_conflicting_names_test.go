package reporter

import (
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/report"
)

func Test_join_conflicting_names(t *testing.T) {
	var (
		ctx, h, s = setup(t)
		m, _, dd  = loadScenarioOwnDM(ctx, s, t, h)
		ff        []*report.Frame
		def       = dd[0]
	)

	ff = loadNoErr(ctx, h, m, def)
	r := ff[0]
	h.a.Equal(5, r.Size())
}
