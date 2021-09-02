package reporter

import (
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/report"
)

func Test1001_loading_basic(t *testing.T) {
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

func Test1002_loading_filtering(t *testing.T) {
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

func Test1003_loading_sorting(t *testing.T) {
	var (
		ctx, h, s = setup(t)
		m, _, dd  = loadScenario(ctx, s, t, h)
		ff        = loadNoErr(ctx, h, m, dd...)
	)

	h.a.Len(ff, 1)
	f := ff[0]
	h.a.Equal(12, f.Size())

	h.a.Equal("id<Record>, first_name<String>, last_name<String>, number_of_numbers<Number>", f.Columns.String())
	checkRows(h, f,
		", Engel, Loritz, 46",
		", Engel, Kiefer, 97",
		", Engel, Kempf, 36",
		", Manu, Specht, 61",
		", Maria, Spannagel, 23",
		", Maria, Königsmann, 61",
		", Maria, Krüger, 99",
		", Sascha, Jans, 38",
		", Sigi, Goldschmidt, 67",
		", Ulli, Haupt, 21",
		", Ulli, Förstner, 87",
		", Ulli, Böhler, 14")

	h.a.Equal("first_name, last_name DESC, id", f.Sort.String())
}

func Test1004_loading_paging(t *testing.T) {
	var (
		ctx, h, s = setup(t)
		m, _, dd  = loadScenario(ctx, s, t, h)
		ff        []*report.Frame
		f         *report.Frame
		def       = dd[0]
	)

	// ^ going up ^
	ff = loadNoErr(ctx, h, m, def)
	h.a.Len(ff, 1)
	f = ff[0]
	h.a.NotNil(f.Paging)
	h.a.NotNil(f.Paging.NextPage)
	h.a.Nil(f.Paging.PrevPage)
	h.a.Equal(5, f.Size())
	checkRows(h, f,
		", Engel_Kempf, Engel",
		", Engel_Kiefer, Engel",
		", Engel_Loritz, Engel",
		", Manu_Specht, Manu",
		", Maria_Krüger, Maria")

	def.Paging.PageCursor = f.Paging.NextPage
	ff = loadNoErr(ctx, h, m, def)
	h.a.Len(ff, 1)
	f = ff[0]
	h.a.NotNil(f.Paging)
	h.a.NotNil(f.Paging.NextPage)
	h.a.NotNil(f.Paging.PrevPage)
	h.a.Equal(5, f.Size())
	checkRows(h, f,
		", Maria_Königsmann, Maria",
		", Maria_Spannagel, Maria",
		", Sascha_Jans, Sascha",
		", Sigi_Goldschmidt, Sigi",
		", Ulli_Böhler, Ulli")

	def.Paging.PageCursor = f.Paging.NextPage
	ff = loadNoErr(ctx, h, m, def)
	h.a.Len(ff, 1)
	f = ff[0]
	h.a.NotNil(f.Paging)
	h.a.Nil(f.Paging.NextPage)
	h.a.NotNil(f.Paging.PrevPage)
	h.a.Equal(2, f.Size())
	checkRows(h, f,
		", Ulli_Förstner, Ulli",
		", Ulli_Haupt, Ulli")

	// v going down v
	def.Paging.PageCursor = f.Paging.PrevPage
	ff = loadNoErr(ctx, h, m, def)
	h.a.Len(ff, 1)
	f = ff[0]
	h.a.NotNil(f.Paging)
	h.a.NotNil(f.Paging.NextPage)
	h.a.NotNil(f.Paging.PrevPage)
	h.a.Equal(5, f.Size())
	checkRows(h, f,
		", Maria_Königsmann, Maria",
		", Maria_Spannagel, Maria",
		", Sascha_Jans, Sascha",
		", Sigi_Goldschmidt, Sigi",
		", Ulli_Böhler, Ulli")

	def.Paging.PageCursor = f.Paging.PrevPage
	ff = loadNoErr(ctx, h, m, def)
	h.a.Len(ff, 1)
	f = ff[0]
	h.a.NotNil(f.Paging)
	h.a.NotNil(f.Paging.NextPage)
	h.a.Nil(f.Paging.PrevPage)
	h.a.Equal(5, f.Size())
	checkRows(h, f,
		", Engel_Kempf, Engel",
		", Engel_Kiefer, Engel",
		", Engel_Loritz, Engel",
		", Manu_Specht, Manu",
		", Maria_Krüger, Maria")
}

func Test1005_loading_paging_no_sort(t *testing.T) {
	var (
		ctx, h, s = setup(t)
		m, _, dd  = loadScenario(ctx, s, t, h)
		ff        []*report.Frame
		f         *report.Frame
		def       = dd[0]
	)

	// ^ going up ^
	ff = loadNoErr(ctx, h, m, def)
	h.a.Len(ff, 1)
	f = ff[0]
	h.a.NotNil(f.Paging)
	h.a.NotNil(f.Paging.NextPage)
	h.a.Nil(f.Paging.PrevPage)
	h.a.Equal(5, f.Size())
	checkRows(h, f,
		", Maria_Königsmann, Maria",
		", Ulli_Haupt, Ulli",
		", Engel_Loritz, Engel",
		", Sascha_Jans, Sascha",
		", Ulli_Böhler, Ulli")

	def.Paging.PageCursor = f.Paging.NextPage
	ff = loadNoErr(ctx, h, m, def)
	h.a.Len(ff, 1)
	f = ff[0]
	h.a.NotNil(f.Paging)
	h.a.NotNil(f.Paging.NextPage)
	h.a.NotNil(f.Paging.PrevPage)
	h.a.Equal(5, f.Size())
	checkRows(h, f,
		", Maria_Spannagel, Maria",
		", Sigi_Goldschmidt, Sigi",
		", Engel_Kempf, Engel",
		", Maria_Krüger, Maria",
		", Manu_Specht, Manu")

	def.Paging.PageCursor = f.Paging.NextPage
	ff = loadNoErr(ctx, h, m, def)
	h.a.Len(ff, 1)
	f = ff[0]
	h.a.NotNil(f.Paging)
	h.a.Nil(f.Paging.NextPage)
	h.a.NotNil(f.Paging.PrevPage)
	h.a.Equal(2, f.Size())
	checkRows(h, f,
		", Ulli_Förstner, Ulli",
		", Engel_Kiefer, Engel")

	// v going down v
	def.Paging.PageCursor = f.Paging.PrevPage
	ff = loadNoErr(ctx, h, m, def)
	h.a.Len(ff, 1)
	f = ff[0]
	h.a.NotNil(f.Paging)
	h.a.NotNil(f.Paging.NextPage)
	h.a.NotNil(f.Paging.PrevPage)
	h.a.Equal(5, f.Size())
	checkRows(h, f,
		", Maria_Spannagel, Maria",
		", Sigi_Goldschmidt, Sigi",
		", Engel_Kempf, Engel",
		", Maria_Krüger, Maria",
		", Manu_Specht, Manu")

	def.Paging.PageCursor = f.Paging.PrevPage
	ff = loadNoErr(ctx, h, m, def)
	h.a.Len(ff, 1)
	f = ff[0]
	h.a.NotNil(f.Paging)
	h.a.NotNil(f.Paging.NextPage)
	h.a.Nil(f.Paging.PrevPage)
	h.a.Equal(5, f.Size())
	checkRows(h, f,
		", Maria_Königsmann, Maria",
		", Ulli_Haupt, Ulli",
		", Engel_Loritz, Engel",
		", Sascha_Jans, Sascha",
		", Ulli_Böhler, Ulli")
}
