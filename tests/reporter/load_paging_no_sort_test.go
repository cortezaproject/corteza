package reporter

import (
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/report"
)

func Test_load_paging_no_sort(t *testing.T) {
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
