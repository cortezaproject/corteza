package reporter

import (
	"testing"
)

func Test_join_paging_nested(t *testing.T) {

	t.Skip("@todo temporarily disabled")

	// var (
	// 	ctx, h, s      = setup(t)
	// 	m, _, dd       = loadScenario(ctx, s, t, h)
	// 	ff             []*report.Frame
	// 	def            = dd[0]
	// 	local, foreign *report.Frame
	// )

	// // // // PAGE 1
	// ff = loadNoErr(ctx, h, m, def)
	// h.a.Len(ff, 5)
	// local = ff[0]
	// ix := indexJoinedResult(ff)
	// _ = ix

	// // local
	// h.a.Equal(2, local.Size())
	// h.a.NotNil(local.Paging)
	// h.a.NotNil(local.Paging.NextPage)
	// checkRows(h, local,
	// 	", aa_01, aa :: 01",
	// 	", aa_02, aa :: 02")

	// foreign = ix["bb/aa/aa_01"]
	// h.a.NotNil(foreign)

	// foreign = ix["bb/aa/aa_02"]
	// h.a.NotNil(foreign)

	// foreign = ix["cc/aa/aa_01"]
	// h.a.NotNil(foreign)

	// foreign = ix["cc/aa/aa_02"]
	// h.a.NotNil(foreign)

	// // // // PAGE 2
	// def.Paging.PageCursor = local.Paging.NextPage
	// ff = loadNoErr(ctx, h, m, def)
	// h.a.Len(ff, 5)
	// local = ff[0]
	// ix = indexJoinedResult(ff)
	// _ = ix

	// // local
	// h.a.Equal(2, local.Size())
	// h.a.NotNil(local.Paging)
	// h.a.NotNil(local.Paging.NextPage)
	// checkRows(h, local,
	// 	", aa_03, aa :: 03",
	// 	", aa_04, aa :: 04")

	// foreign = ix["bb/aa/aa_03"]
	// h.a.NotNil(foreign)

	// foreign = ix["bb/aa/aa_04"]
	// h.a.NotNil(foreign)

	// foreign = ix["cc/aa/aa_03"]
	// h.a.NotNil(foreign)

	// foreign = ix["cc/aa/aa_04"]
	// h.a.NotNil(foreign)

	// // // // PAGE 3
	// def.Paging.PageCursor = local.Paging.NextPage
	// ff = loadNoErr(ctx, h, m, def)
	// h.a.Len(ff, 3)
	// local = ff[0]
	// ix = indexJoinedResult(ff)
	// _ = ix

	// // local
	// h.a.Equal(1, local.Size())
	// h.a.Nil(local.Paging)
	// checkRows(h, local,
	// 	", aa_05, aa :: 05")

	// foreign = ix["bb/aa/aa_05"]
	// h.a.NotNil(foreign)

	// foreign = ix["cc/aa/aa_05"]
	// h.a.NotNil(foreign)
}
