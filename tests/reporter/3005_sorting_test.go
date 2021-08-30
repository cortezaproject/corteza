package reporter

import (
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/report"
)

func Test3005_sorting(t *testing.T) {
	var (
		ctx, h, s      = setup(t)
		m, _, dd       = loadScenario(ctx, s, t, h)
		ff             []*report.Frame
		local, foreign *report.Frame
	)

	ff = loadNoErr(ctx, h, m, dd...)
	h.a.Len(ff, 7)
	local = ff[0]
	ix := indexJoinedResult(ff)
	_ = ix

	// local
	h.a.Equal(12, local.Size())
	h.a.Equal("id<Record>, join_key<String>, first_name<String>, last_name<String>", local.Columns.String())
	h.a.Equal("first_name, last_name DESC", local.Sort.String())
	checkRows(h, local,
		", Engel, Loritz",
		", Engel, Kiefer",
		", Engel, Kempf",
		", Manu, Specht",
		", Maria, Spannagel",
		", Maria, Königsmann",
		", Maria, Krüger",
		", Sascha, Jans",
		", Sigi, Goldschmidt",
		", Ulli, Haupt",
		", Ulli, Förstner",
		", Ulli, Böhler")

	// Maria_Königsmann
	foreign = ix["Maria_Königsmann"]
	h.a.NotNil(foreign)
	h.a.Equal("id<Record>, usr<String>, name<String>, type<Select>, cost<Number>, time_spent<Number>", foreign.Columns.String())
	h.a.Equal("type", foreign.Sort.String())
	checkRows(h, foreign,
		", Maria_Königsmann, u1 j1, a, 10, 2",
		", Maria_Königsmann, u1 j5, a, 4, 4",
		", Maria_Königsmann, u1 j2, b, 20, 5",
		", Maria_Königsmann, u1 j6, b, 25, 25",
		", Maria_Königsmann, u1 j7, b, 9, 91",
		", Maria_Königsmann, u1 j4, c, 3, 0",
		", Maria_Königsmann, u1 j3, d, 11, 1")

	// Engel_Loritz
	foreign = ix["Engel_Loritz"]
	h.a.NotNil(foreign)
	h.a.Equal("id<Record>, usr<String>, name<String>, type<Select>, cost<Number>, time_spent<Number>", foreign.Columns.String())
	h.a.Equal("type", foreign.Sort.String())
	checkRows(h, foreign,
		", Engel_Loritz, u3 j1, a, 10, 1",
		", Engel_Loritz, u3 j2, a, 0, 0",
		", Engel_Loritz, u3 j3, a, 19, 99")

	// Sigi_Goldschmidt
	foreign = ix["Sigi_Goldschmidt"]
	h.a.NotNil(foreign)
	h.a.Equal("id<Record>, usr<String>, name<String>, type<Select>, cost<Number>, time_spent<Number>", foreign.Columns.String())
	h.a.Equal("type", foreign.Sort.String())
	checkRows(h, foreign,
		", Sigi_Goldschmidt, u7 j2, a, 10, 21",
		", Sigi_Goldschmidt, u7 j3, b, 10, 99",
		", Sigi_Goldschmidt, u7 j1, d, 10, 29")

	// Engel_Kiefer
	foreign = ix["Engel_Kiefer"]
	h.a.NotNil(foreign)
	h.a.Equal("id<Record>, usr<String>, name<String>, type<Select>, cost<Number>, time_spent<Number>", foreign.Columns.String())
	h.a.Equal("type", foreign.Sort.String())
	checkRows(h, foreign,
		", Engel_Kiefer, u12 j1, a, 42, 69",
		", Engel_Kiefer, u12 j4, a, 35, 26",
		", Engel_Kiefer, u12 j5, a, 34, 29",
		", Engel_Kiefer, u12 j9, a, 71, 90",
		", Engel_Kiefer, u12 j2, b, 65, 99",
		", Engel_Kiefer, u12 j8, b, 74, 39",
		", Engel_Kiefer, u12 j3, c, 38, 71",
		", Engel_Kiefer, u12 j10, c, 79, 25",
		", Engel_Kiefer, u12 j11, c, 19, 66",
		", Engel_Kiefer, u12 j6, d, 92, 16",
		", Engel_Kiefer, u12 j7, d, 14, 71")

	// Manu_Specht
	foreign = ix["Manu_Specht"]
	h.a.NotNil(foreign)
	h.a.Equal("id<Record>, usr<String>, name<String>, type<Select>, cost<Number>, time_spent<Number>", foreign.Columns.String())
	h.a.Equal("type", foreign.Sort.String())
	checkRows(h, foreign,
		", Manu_Specht, u10 j3, b, 53, 12",
		", Manu_Specht, u10 j4, b, 60, 22",
		", Manu_Specht, u10 j2, c, 83, 70",
		", Manu_Specht, u10 j1, d, 45, 56")

	// Ulli_Böhler
	foreign = ix["Ulli_Böhler"]
	h.a.NotNil(foreign)
	h.a.Equal("id<Record>, usr<String>, name<String>, type<Select>, cost<Number>, time_spent<Number>", foreign.Columns.String())
	h.a.Equal("type", foreign.Sort.String())
	checkRows(h, foreign,
		", Ulli_Böhler, u5 j1, a, 1, 2")
}
