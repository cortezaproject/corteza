package reporter

import (
	"testing"
)

func Test_load_sorting(t *testing.T) {
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
