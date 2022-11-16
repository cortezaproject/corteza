package reporter

import (
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/report"
)

func Test_join_dup_keys(t *testing.T) {
	var (
		ctx, h, s = setup(t)
		m, _, dd  = loadScenarioOwnDM(ctx, s, t, h)
		ff        []*report.Frame

		local *report.Frame
		def   = dd[0]
	)

	// // // PAGE 1
	ff = loadNoErr(ctx, h, m, def)
	h.a.Len(ff, 3)
	local = ff[0]

	// local
	h.a.Equal(2, local.Size())
	h.a.NotNil(local.Paging)
	h.a.NotNil(local.Paging.NextPage)
	checkRows(h, local,
		"a, 11",
		"a, 12")

	h.a.Equal("a", ff[1].RefValue)
	h.a.Equal("a", ff[2].RefValue)

	// // // PAGE 2
	def.Paging.PageCursor = local.Paging.NextPage
	ff = loadNoErr(ctx, h, m, def)
	h.a.Len(ff, 3)
	local = ff[0]

	// local
	h.a.Equal(2, local.Size())
	h.a.NotNil(local.Paging)
	h.a.NotNil(local.Paging.NextPage)
	checkRows(h, local,
		"a, 13",
		"b, 14")

	h.a.Equal("a", ff[1].RefValue)
	h.a.Equal("b", ff[2].RefValue)

	// // // PAGE 1
	def.Paging.PageCursor = local.Paging.NextPage
	ff = loadNoErr(ctx, h, m, def)
	h.a.Len(ff, 2)
	local = ff[0]

	// local
	h.a.Equal(1, local.Size())
	h.a.Nil(local.Paging)
	checkRows(h, local,
		"c, 15")

	h.a.Equal("c", ff[1].RefValue)
}
