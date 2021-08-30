package reporter

import (
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/report"
)

func Test3004_joining_paging(t *testing.T) {
	var (
		ctx, h, s      = setup(t)
		m, _, dd       = loadScenario(ctx, s, t, h)
		ff             []*report.Frame
		def            = dd[0]
		local, foreign *report.Frame
	)

	// // // PAGE 1
	ff = loadNoErr(ctx, h, m, def)
	h.a.Len(ff, 3)
	local = ff[0]
	ix := indexJoinedResult(ff)
	_ = ix

	// local
	h.a.Equal(2, local.Size())
	h.a.NotNil(local.Paging)
	h.a.NotNil(local.Paging.NextPage)
	checkRows(h, local,
		", aa_05, aa :: 05",
		", aa_04, aa :: 04")

	foreign = ix["aa_05"]
	h.a.NotNil(foreign)

	foreign = ix["aa_04"]
	h.a.NotNil(foreign)

	// // // PAGE 2
	def.Paging.PageCursor = local.Paging.NextPage
	ff = loadNoErr(ctx, h, m, def)
	h.a.Len(ff, 3)
	local = ff[0]
	ix = indexJoinedResult(ff)
	_ = ix

	// local
	h.a.Equal(2, local.Size())
	h.a.NotNil(local.Paging)
	h.a.NotNil(local.Paging.NextPage)
	checkRows(h, local,
		", aa_03, aa :: 03",
		", aa_02, aa :: 02")

	foreign = ix["aa_03"]
	h.a.NotNil(foreign)

	foreign = ix["aa_02"]
	h.a.NotNil(foreign)

	// // // PAGE 1
	def.Paging.PageCursor = local.Paging.NextPage
	ff = loadNoErr(ctx, h, m, def)
	h.a.Len(ff, 2)
	local = ff[0]
	ix = indexJoinedResult(ff)
	_ = ix

	// local
	h.a.Equal(1, local.Size())
	h.a.Nil(local.Paging)
	checkRows(h, local,
		", aa_01, aa :: 01")

	foreign = ix["aa_01"]
	h.a.NotNil(foreign)
}
