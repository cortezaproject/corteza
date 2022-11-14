package reporter

import (
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/report"
)

func Test_load_paging(t *testing.T) {
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
