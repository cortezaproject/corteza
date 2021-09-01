package reporter

import (
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/report"
)

func Test3002_joining_base_nested(t *testing.T) {

	t.Skip("@todo temporarily disabled")

	var (
		ctx, h, s      = setup(t)
		m, _, dd       = loadScenario(ctx, s, t, h)
		ff             = loadNoErr(ctx, h, m, dd...)
		local, foreign *report.Frame
	)

	// The joining here looks like this:
	//
	//            (nested)
	//  (nested_aux)    (cc)
	// (aa)     (bb)

	h.a.Len(ff, 11)

	ix := indexJoinedResult(ff)
	_ = ix

	// // joined -- the initial join

	// local
	local = ff[0]
	h.a.Equal("pk<String>, label<String>", local.Columns.OmitSys().String())
	h.a.Equal("joined", local.Source)
	h.a.Equal("aa", local.Ref)
	checkRows(h, local,
		"aa_01, aa :: 01",
		"aa_02, aa :: 02",
		"aa_03, aa :: 03",
		"aa_04, aa :: 04",
		"aa_05, aa :: 05")

	// aa_01
	foreign = ix["bb/aa/aa_01"]
	h.a.NotNil(foreign)
	h.a.Equal("pk<String>, fk_a<String>, label<String>", foreign.Columns.OmitSys().String())
	h.a.Equal("joined", foreign.Source)
	h.a.Equal("bb", foreign.Ref)
	h.a.Equal("pk", foreign.RelColumn)
	h.a.Equal("aa_01", foreign.RefValue)
	checkRows(h, foreign,
		"bb_01, aa_01, bb :: 01",
		"bb_02, aa_01, bb :: 02",
		"bb_03, aa_01, bb :: 03")

	// aa_02
	foreign = ix["bb/aa/aa_02"]
	h.a.NotNil(foreign)
	h.a.Equal("pk<String>, fk_a<String>, label<String>", foreign.Columns.OmitSys().String())
	h.a.Equal("joined", foreign.Source)
	h.a.Equal("bb", foreign.Ref)
	h.a.Equal("pk", foreign.RelColumn)
	h.a.Equal("aa_02", foreign.RefValue)
	checkRows(h, foreign,
		"bb_04, aa_02, bb :: 04",
		"bb_05, aa_02, bb :: 05")

	// aa_03
	foreign = ix["bb/aa/aa_03"]
	h.a.NotNil(foreign)
	h.a.Equal("pk<String>, fk_a<String>, label<String>", foreign.Columns.OmitSys().String())
	h.a.Equal("joined", foreign.Source)
	h.a.Equal("bb", foreign.Ref)
	h.a.Equal("pk", foreign.RelColumn)
	h.a.Equal("aa_03", foreign.RefValue)
	checkRows(h, foreign,
		"bb_06, aa_03, bb :: 06")

	// aa_04
	foreign = ix["bb/aa/aa_04"]
	h.a.NotNil(foreign)
	h.a.Equal("pk<String>, fk_a<String>, label<String>", foreign.Columns.OmitSys().String())
	h.a.Equal("joined", foreign.Source)
	h.a.Equal("bb", foreign.Ref)
	h.a.Equal("pk", foreign.RelColumn)
	h.a.Equal("aa_04", foreign.RefValue)
	checkRows(h, foreign,
		"bb_07, aa_04, bb :: 07")

	// aa_05
	foreign = ix["bb/aa/aa_05"]
	h.a.NotNil(foreign)
	h.a.Equal("pk<String>, fk_a<String>, label<String>", foreign.Columns.OmitSys().String())
	h.a.Equal("joined", foreign.Source)
	h.a.Equal("bb", foreign.Ref)
	h.a.Equal("pk", foreign.RelColumn)
	h.a.Equal("aa_05", foreign.RefValue)
	checkRows(h, foreign,
		"bb_08, aa_05, bb :: 08",
		"bb_09, aa_05, bb :: 09",
		"bb_10, aa_05, bb :: 10")

	// The other foreign

	// aa_01
	foreign = ix["cc/aa/aa_01"]
	h.a.NotNil(foreign)
	h.a.Equal("pk<String>, fk_a<String>, label<String>", foreign.Columns.OmitSys().String())
	h.a.Equal("joined", foreign.Source)
	h.a.Equal("cc", foreign.Ref)
	h.a.Equal("pk", foreign.RelColumn)
	h.a.Equal("aa_01", foreign.RefValue)
	checkRows(h, foreign,
		"cc_01, aa_01, cc :: 01",
		"cc_02, aa_01, cc :: 02")

	// aa_02
	foreign = ix["cc/aa/aa_02"]
	h.a.NotNil(foreign)
	h.a.Equal("pk<String>, fk_a<String>, label<String>", foreign.Columns.OmitSys().String())
	h.a.Equal("joined", foreign.Source)
	h.a.Equal("cc", foreign.Ref)
	h.a.Equal("pk", foreign.RelColumn)
	h.a.Equal("aa_02", foreign.RefValue)
	checkRows(h, foreign,
		"cc_03, aa_02, cc :: 03",
		"cc_04, aa_02, cc :: 04")

	// aa_03
	foreign = ix["cc/aa/aa_03"]
	h.a.NotNil(foreign)
	h.a.Equal("pk<String>, fk_a<String>, label<String>", foreign.Columns.OmitSys().String())
	h.a.Equal("joined", foreign.Source)
	h.a.Equal("cc", foreign.Ref)
	h.a.Equal("pk", foreign.RelColumn)
	h.a.Equal("aa_03", foreign.RefValue)
	checkRows(h, foreign,
		"cc_05, aa_03, cc :: 05",
		"cc_06, aa_03, cc :: 06")

	// aa_04
	foreign = ix["cc/aa/aa_04"]
	h.a.NotNil(foreign)
	h.a.Equal("pk<String>, fk_a<String>, label<String>", foreign.Columns.OmitSys().String())
	h.a.Equal("joined", foreign.Source)
	h.a.Equal("cc", foreign.Ref)
	h.a.Equal("pk", foreign.RelColumn)
	h.a.Equal("aa_04", foreign.RefValue)
	checkRows(h, foreign,
		"cc_07, aa_04, cc :: 07")

	// aa_05
	foreign = ix["cc/aa/aa_05"]
	h.a.NotNil(foreign)
	h.a.Equal("pk<String>, fk_a<String>, label<String>", foreign.Columns.OmitSys().String())
	h.a.Equal("joined", foreign.Source)
	h.a.Equal("cc", foreign.Ref)
	h.a.Equal("pk", foreign.RelColumn)
	h.a.Equal("aa_05", foreign.RefValue)
	checkRows(h, foreign,
		"cc_08, aa_05, cc :: 08")
}
