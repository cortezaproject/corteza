package reporter

import (
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/report"
)

func Test3001_joining_basic(t *testing.T) {
	var (
		ctx, h, s      = setup(t)
		m, _, dd       = loadScenario(ctx, s, t, h)
		ff             = loadNoErr(ctx, h, m, dd...)
		local, foreign *report.Frame
	)

	h.a.Len(ff, 7)
	local = ff[0]
	ix := indexJoinedResult(ff)
	_ = ix

	// local
	h.a.Equal(12, local.Size())
	h.a.Equal("id<Record>, join_key<String>, first_name<String>, last_name<String>", local.Columns.String())
	checkRows(h, local,
		", Engel_Loritz, Engel, Loritz",
		", Engel_Kempf, Engel, Kempf",
		", Engel_Kiefer, Engel, Kiefer",
		", Manu_Specht, Manu, Specht",
		", Maria_Königsmann, Maria, Königsmann",
		", Maria_Spannagel, Maria, Spannagel",
		", Maria_Krüger, Maria, Krüger",
		", Sascha_Jans, Sascha, Jans",
		", Sigi_Goldschmidt, Sigi, Goldschmidt",
		", Ulli_Haupt, Ulli, Haupt",
		", Ulli_Böhler, Ulli, Böhler",
		", Ulli_Förstner, Ulli, Förstner")

	// Maria_Königsmann
	foreign = ix["Maria_Königsmann"]
	h.a.NotNil(foreign)
	h.a.Equal("id<Record>, usr<String>, name<String>, type<Select>, cost<Number>, time_spent<Number>", foreign.Columns.String())
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
	checkRows(h, foreign,
		", Engel_Loritz, u3 j1, a, 10, 1",
		", Engel_Loritz, u3 j2, a, 0, 0",
		", Engel_Loritz, u3 j3, a, 19, 99")

	// Sigi_Goldschmidt
	foreign = ix["Sigi_Goldschmidt"]
	h.a.NotNil(foreign)
	h.a.Equal("id<Record>, usr<String>, name<String>, type<Select>, cost<Number>, time_spent<Number>", foreign.Columns.String())
	checkRows(h, foreign,
		", Sigi_Goldschmidt, u7 j2, a, 10, 21",
		", Sigi_Goldschmidt, u7 j3, b, 10, 99",
		", Sigi_Goldschmidt, u7 j1, d, 10, 29")

	// Engel_Kiefer
	foreign = ix["Engel_Kiefer"]
	h.a.NotNil(foreign)
	h.a.Equal("id<Record>, usr<String>, name<String>, type<Select>, cost<Number>, time_spent<Number>", foreign.Columns.String())
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
	checkRows(h, foreign,
		", Manu_Specht, u10 j3, b, 53, 12",
		", Manu_Specht, u10 j4, b, 60, 22",
		", Manu_Specht, u10 j2, c, 83, 70",
		", Manu_Specht, u10 j1, d, 45, 56")

	// Ulli_Böhler
	foreign = ix["Ulli_Böhler"]
	h.a.NotNil(foreign)
	h.a.Equal("id<Record>, usr<String>, name<String>, type<Select>, cost<Number>, time_spent<Number>", foreign.Columns.String())
	checkRows(h, foreign,
		", Ulli_Böhler, u5 j1, a, 1, 2")
}

func Test3002_joining_filtering(t *testing.T) {
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
	foreign = ix["Engel_Kiefer"]
	h.a.NotNil(foreign)
	h.a.Equal(4, foreign.Size())
	h.a.Equal("id<Record>, usr<String>, name<String>, type<Select>, cost<Number>, time_spent<Number>", foreign.Columns.String())
	checkRows(h, foreign,
		", Engel_Kiefer, u12 j1, a, 42, 69",
		", Engel_Kiefer, u12 j4, a, 35, 26",
		", Engel_Kiefer, u12 j5, a, 34, 29",
		", Engel_Kiefer, u12 j9, a, 71, 90")

	// Engel_Loritz
	foreign = ix["Maria_Königsmann"]
	h.a.NotNil(foreign)
	h.a.Equal(2, foreign.Size())
	h.a.Equal("id<Record>, usr<String>, name<String>, type<Select>, cost<Number>, time_spent<Number>", foreign.Columns.String())
	checkRows(h, foreign,
		", Maria_Königsmann, u1 j1, a, 10, 2",
		", Maria_Königsmann, u1 j5, a, 4, 4")

}
func Test3003_joining_sorting_local(t *testing.T) {
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

func Test3004_joining_sorting_foreign(t *testing.T) {
	var (
		ctx, h, s      = setup(t)
		m, _, dd       = loadScenario(ctx, s, t, h)
		ff             []*report.Frame
		local, foreign *report.Frame
	)

	// @todo !!! for now local DS sorting based on the foreign DS includes
	//       only the rows that have a match.
	//       This will do for now, but should be improved in the future.

	ff = loadNoErr(ctx, h, m, dd...)
	h.a.Len(ff, 7)
	local = ff[0]
	ix := indexJoinedResult(ff)
	_ = ix

	// local
	h.a.Equal(6, local.Size())
	h.a.Equal("id<Record>, join_key<String>, first_name<String>, last_name<String>", local.Columns.String())
	// h.a.Equal("jobs.type, first_name, last_name DESC", local.Sort.String())
	checkRows(h, local,
		", Engel, Loritz",
		", Engel, Kiefer",
		", Maria, Königsmann",
		", Sigi, Goldschmidt",
		", Ulli, Böhler",
		", Manu, Specht")

	// Maria_Königsmann
	foreign = ix["Maria_Königsmann"]
	h.a.NotNil(foreign)
	h.a.Equal("id<Record>, usr<String>, name<String>, type<Select>, cost<Number>, time_spent<Number>", foreign.Columns.String())
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
	checkRows(h, foreign,
		", Engel_Loritz, u3 j1, a, 10, 1",
		", Engel_Loritz, u3 j2, a, 0, 0",
		", Engel_Loritz, u3 j3, a, 19, 99")

	// Sigi_Goldschmidt
	foreign = ix["Sigi_Goldschmidt"]
	h.a.NotNil(foreign)
	h.a.Equal("id<Record>, usr<String>, name<String>, type<Select>, cost<Number>, time_spent<Number>", foreign.Columns.String())
	checkRows(h, foreign,
		", Sigi_Goldschmidt, u7 j2, a, 10, 21",
		", Sigi_Goldschmidt, u7 j3, b, 10, 99",
		", Sigi_Goldschmidt, u7 j1, d, 10, 29")

	// Engel_Kiefer
	foreign = ix["Engel_Kiefer"]
	h.a.NotNil(foreign)
	h.a.Equal("id<Record>, usr<String>, name<String>, type<Select>, cost<Number>, time_spent<Number>", foreign.Columns.String())
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
	checkRows(h, foreign,
		", Manu_Specht, u10 j3, b, 53, 12",
		", Manu_Specht, u10 j4, b, 60, 22",
		", Manu_Specht, u10 j2, c, 83, 70",
		", Manu_Specht, u10 j1, d, 45, 56")

	// Ulli_Böhler
	foreign = ix["Ulli_Böhler"]
	h.a.NotNil(foreign)
	h.a.Equal("id<Record>, usr<String>, name<String>, type<Select>, cost<Number>, time_spent<Number>", foreign.Columns.String())
	checkRows(h, foreign,
		", Ulli_Böhler, u5 j1, a, 1, 2")
}

func Test3005_joining_paging(t *testing.T) {
	var (
		ctx, h, s      = setup(t)
		m, _, dd       = loadScenario(ctx, s, t, h)
		ff             []*report.Frame
		local, foreign *report.Frame
	)

	// ^ going up ^

	// // // PAGE 1
	ff = loadNoErr(ctx, h, m, dd...)
	h.a.Len(ff, 4)
	local = ff[0]
	ix := indexJoinedResult(ff)
	_ = ix

	// local
	h.a.Equal(5, local.Size())
	h.a.NotNil(local.Paging)
	h.a.NotNil(local.Paging.NextPage)
	h.a.Nil(local.Paging.PrevPage)
	checkRows(h, local,
		", Engel_Kempf, Engel, Kempf",
		", Engel_Kiefer, Engel, Kiefer",
		", Engel_Loritz, Engel, Loritz",
		", Manu_Specht, Manu, Specht",
		", Maria_Krüger, Maria, Krüger")

	// Manu_Specht
	foreign = ix["Manu_Specht"]
	h.a.NotNil(foreign)
	h.a.NotNil(foreign.Paging)
	h.a.NotNil(foreign.Paging.NextPage)
	h.a.Nil(foreign.Paging.PrevPage)
	checkRows(h, foreign,
		", Manu_Specht, u10 j1, d, 45, 56",
		", Manu_Specht, u10 j2, c, 83, 70")

	// Engel_Kiefer
	foreign = ix["Engel_Kiefer"]
	h.a.NotNil(foreign)
	h.a.NotNil(foreign.Paging)
	h.a.NotNil(foreign.Paging.NextPage)
	h.a.Nil(foreign.Paging.PrevPage)
	checkRows(h, foreign,
		", Engel_Kiefer, u12 j1, a, 42, 69",
		", Engel_Kiefer, u12 j10, c, 79, 25")

	// Engel_Loritz
	foreign = ix["Engel_Loritz"]
	h.a.NotNil(foreign)
	h.a.NotNil(foreign.Paging)
	h.a.NotNil(foreign.Paging.NextPage)
	h.a.Nil(foreign.Paging.PrevPage)
	checkRows(h, foreign,
		", Engel_Loritz, u3 j1, a, 10, 1",
		", Engel_Loritz, u3 j2, a, 0, 0")

	// // // PAGE 2
	dd[0].Paging.PageCursor = local.Paging.NextPage
	ff = loadNoErr(ctx, h, m, dd...)

	h.a.Len(ff, 4)
	local = ff[0]
	ix = indexJoinedResult(ff)
	_ = ix

	// local
	h.a.Equal(5, local.Size())
	h.a.NotNil(local.Paging)
	h.a.NotNil(local.Paging.NextPage)
	h.a.NotNil(local.Paging.PrevPage)
	checkRows(h, local,
		", Maria_Königsmann, Maria, Königsmann",
		", Maria_Spannagel, Maria, Spannagel",
		", Sascha_Jans, Sascha, Jans",
		", Sigi_Goldschmidt, Sigi, Goldschmidt",
		", Ulli_Böhler, Ulli, Böhler",
		", Ulli_Förstner, Ulli, Förstner")

	// Maria_Königsmann
	foreign = ix["Maria_Königsmann"]
	h.a.NotNil(foreign)
	h.a.NotNil(foreign.Paging)
	h.a.NotNil(foreign.Paging.NextPage)
	h.a.Nil(foreign.Paging.PrevPage)
	checkRows(h, foreign,
		", Maria_Königsmann, u1 j1, a, 10, 2",
		", Maria_Königsmann, u1 j2, b, 20, 5")

	// Engel_Kiefer
	foreign = ix["Ulli_Böhler"]
	h.a.NotNil(foreign)
	h.a.Nil(foreign.Paging)
	checkRows(h, foreign,
		", Ulli_Böhler, u5 j1, a, 1, 2")

	// Sigi_Goldschmidt
	foreign = ix["Sigi_Goldschmidt"]
	h.a.NotNil(foreign)
	h.a.NotNil(foreign.Paging)
	h.a.NotNil(foreign.Paging.NextPage)
	h.a.Nil(foreign.Paging.PrevPage)
	checkRows(h, foreign,
		", Sigi_Goldschmidt, u7 j1, d, 10, 29",
		", Sigi_Goldschmidt, u7 j2, a, 10, 21")

	// // // PAGE 3
	dd[0].Paging.PageCursor = local.Paging.NextPage
	ff = loadNoErr(ctx, h, m, dd...)

	h.a.Len(ff, 1)
	local = ff[0]
	ix = indexJoinedResult(ff)
	_ = ix

	// local
	h.a.Equal(2, local.Size())
	h.a.NotNil(local.Paging)
	h.a.Nil(local.Paging.NextPage)
	h.a.NotNil(local.Paging.PrevPage)
	checkRows(h, local,
		", Ulli_Förstner, Ulli, Förstner",
		", Ulli_Haupt, Ulli, Haupt")

	// v going down v

	// // // PAGE 2
	dd[0].Paging.PageCursor = local.Paging.PrevPage
	ff = loadNoErr(ctx, h, m, dd...)

	h.a.Len(ff, 4)
	local = ff[0]
	ix = indexJoinedResult(ff)
	_ = ix

	// local
	h.a.Equal(5, local.Size())
	h.a.NotNil(local.Paging)
	h.a.NotNil(local.Paging.NextPage)
	h.a.NotNil(local.Paging.PrevPage)
	checkRows(h, local,
		", Maria_Königsmann, Maria, Königsmann",
		", Maria_Spannagel, Maria, Spannagel",
		", Sascha_Jans, Sascha, Jans",
		", Sigi_Goldschmidt, Sigi, Goldschmidt",
		", Ulli_Böhler, Ulli, Böhler",
		", Ulli_Förstner, Ulli, Förstner")

	// Maria_Königsmann
	foreign = ix["Maria_Königsmann"]
	h.a.NotNil(foreign)
	h.a.NotNil(foreign.Paging)
	h.a.NotNil(foreign.Paging.NextPage)
	h.a.Nil(foreign.Paging.PrevPage)
	checkRows(h, foreign,
		", Maria_Königsmann, u1 j1, a, 10, 2",
		", Maria_Königsmann, u1 j2, b, 20, 5")

	// Engel_Kiefer
	foreign = ix["Ulli_Böhler"]
	h.a.NotNil(foreign)
	h.a.Nil(foreign.Paging)
	checkRows(h, foreign,
		", Ulli_Böhler, u5 j1, a, 1, 2")

	// Sigi_Goldschmidt
	foreign = ix["Sigi_Goldschmidt"]
	h.a.NotNil(foreign)
	h.a.NotNil(foreign.Paging)
	h.a.NotNil(foreign.Paging.NextPage)
	h.a.Nil(foreign.Paging.PrevPage)
	checkRows(h, foreign,
		", Sigi_Goldschmidt, u7 j1, d, 10, 29",
		", Sigi_Goldschmidt, u7 j2, a, 10, 21")

	// // // PAGE 1
	dd[0].Paging.PageCursor = local.Paging.PrevPage
	ff = loadNoErr(ctx, h, m, dd...)
	h.a.Len(ff, 4)
	local = ff[0]
	ix = indexJoinedResult(ff)
	_ = ix

	// local
	h.a.Equal(5, local.Size())
	h.a.NotNil(local.Paging)
	h.a.NotNil(local.Paging.NextPage)
	h.a.Nil(local.Paging.PrevPage)
	checkRows(h, local,
		", Engel_Kempf, Engel, Kempf",
		", Engel_Kiefer, Engel, Kiefer",
		", Engel_Loritz, Engel, Loritz",
		", Manu_Specht, Manu, Specht",
		", Maria_Krüger, Maria, Krüger")

	// Manu_Specht
	foreign = ix["Manu_Specht"]
	h.a.NotNil(foreign)
	h.a.NotNil(foreign.Paging)
	h.a.NotNil(foreign.Paging.NextPage)
	h.a.Nil(foreign.Paging.PrevPage)
	checkRows(h, foreign,
		", Manu_Specht, u10 j1, d, 45, 56",
		", Manu_Specht, u10 j2, c, 83, 70")

	// Engel_Kiefer
	foreign = ix["Engel_Kiefer"]
	h.a.NotNil(foreign)
	h.a.NotNil(foreign.Paging)
	h.a.NotNil(foreign.Paging.NextPage)
	h.a.Nil(foreign.Paging.PrevPage)
	checkRows(h, foreign,
		", Engel_Kiefer, u12 j1, a, 42, 69",
		", Engel_Kiefer, u12 j10, c, 79, 25")

	// Engel_Loritz
	foreign = ix["Engel_Loritz"]
	h.a.NotNil(foreign)
	h.a.NotNil(foreign.Paging)
	h.a.NotNil(foreign.Paging.NextPage)
	h.a.Nil(foreign.Paging.PrevPage)
	checkRows(h, foreign,
		", Engel_Loritz, u3 j1, a, 10, 1",
		", Engel_Loritz, u3 j2, a, 0, 0")
}

func indexJoinedResult(ff []*report.Frame) map[string]*report.Frame {
	out := make(map[string]*report.Frame)
	// the first one is the local ds
	for _, f := range ff[1:] {
		out[f.RefValue] = f
	}

	return out
}
