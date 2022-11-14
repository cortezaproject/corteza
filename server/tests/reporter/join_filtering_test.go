package reporter

import (
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/report"
)

func Test_join_filtering(t *testing.T) {
	var (
		ctx, h, s      = setup(t)
		m, _, dd       = loadScenario(ctx, s, t, h)
		ff             = loadNoErr(ctx, h, m, dd...)
		local, foreign *report.Frame
	)

	h.a.Len(ff, 3)
	local = ff[0]
	ix := indexJoinedResult(ff)
	_ = ix

	// local
	h.a.Equal(2, local.Size())
	h.a.Equal("id<Record>, join_key<String>, first_name<String>, last_name<String>", local.Columns.String())
	checkRows(h, local,
		", Engel_Kiefer, Engel, Kiefer",
		", Maria_Königsmann, Maria, Königsmann")

	// Engel_Kiefer
	foreign = ix["jobs/users/Engel_Kiefer"]
	h.a.NotNil(foreign)
	h.a.Equal(4, foreign.Size())
	h.a.Equal("id<Record>, usr<String>, name<String>, type<Select>, cost<Number>, time_spent<Number>", foreign.Columns.String())
	checkRows(h, foreign,
		", Engel_Kiefer, u12 j1, a, 42, 69",
		", Engel_Kiefer, u12 j4, a, 35, 26",
		", Engel_Kiefer, u12 j5, a, 34, 29",
		", Engel_Kiefer, u12 j9, a, 71, 90")

	// Engel_Loritz
	foreign = ix["jobs/users/Maria_Königsmann"]
	h.a.NotNil(foreign)
	h.a.Equal(2, foreign.Size())
	h.a.Equal("id<Record>, usr<String>, name<String>, type<Select>, cost<Number>, time_spent<Number>", foreign.Columns.String())
	checkRows(h, foreign,
		", Maria_Königsmann, u1 j1, a, 10, 2",
		", Maria_Königsmann, u1 j5, a, 4, 4")
}
