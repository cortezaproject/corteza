package reporter

import (
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/report"
)

func Test9001_filtering_ast(t *testing.T) {
	var (
		ctx, h, s = setup(t)
		m, _, dd  = loadScenario(ctx, s, t, h)
		ff        = loadNoErrMulti(ctx, h, m, dd...)
		f         *report.Frame
	)

	h.a.Len(ff, 1)

	f = ff[0]
	h.a.Len(f.Rows, 5)
	h.a.Equal("first_name<String>, last_name<String>", f.Columns.String())
	checkRows(h, f,
		"Maria, Königsmann",
		"Maria, Spannagel",
		"Engel, Kempf",
		"Maria, Krüger",
		"Engel, Kiefer")
}

func Test9002_filtering_expr(t *testing.T) {
	var (
		ctx, h, s = setup(t)
		m, _, dd  = loadScenario(ctx, s, t, h)
		ff        = loadNoErrMulti(ctx, h, m, dd...)
		f         *report.Frame
	)

	h.a.Len(ff, 1)

	f = ff[0]
	h.a.Len(f.Rows, 5)
	h.a.Equal("first_name<String>, last_name<String>", f.Columns.String())
	checkRows(h, f,
		"Maria, Königsmann",
		"Maria, Spannagel",
		"Engel, Kempf",
		"Maria, Krüger",
		"Engel, Kiefer")
}

func Test9003_filtering_ast_expr(t *testing.T) {
	var (
		ctx, h, s = setup(t)
		m, _, dd  = loadScenario(ctx, s, t, h)
		ff        = loadNoErrMulti(ctx, h, m, dd...)
		f         *report.Frame
	)

	h.a.Len(ff, 1)

	f = ff[0]
	h.a.Len(f.Rows, 5)
	h.a.Equal("first_name<String>, last_name<String>", f.Columns.String())
	checkRows(h, f,
		"Maria, Königsmann",
		"Maria, Spannagel",
		"Engel, Kempf",
		"Maria, Krüger",
		"Engel, Kiefer")
}

func Test9004_filtering_fnc_ast(t *testing.T) {
	var (
		ctx, h, s = setup(t)
		m, _, dd  = loadScenario(ctx, s, t, h)
		ff        = loadNoErrMulti(ctx, h, m, dd...)
		f         *report.Frame
	)

	h.a.Len(ff, 1)

	f = ff[0]
	h.a.Len(f.Rows, 7)
	h.a.Equal("first_name<String>, last_name<String>", f.Columns.String())
	checkRows(h, f,
		"Ulli, Haupt",
		"Engel, Loritz",
		"Maria, Spannagel",
		"Sigi, Goldschmidt",
		"Engel, Kempf",
		"Manu, Specht",
		"Ulli, Förstner")
}

func Test9005_filtering_fnc_expr(t *testing.T) {
	var (
		ctx, h, s = setup(t)
		m, _, dd  = loadScenario(ctx, s, t, h)
		ff        = loadNoErrMulti(ctx, h, m, dd...)
		f         *report.Frame
	)

	h.a.Len(ff, 1)

	f = ff[0]
	h.a.Len(f.Rows, 7)
	h.a.Equal("first_name<String>, last_name<String>", f.Columns.String())
	checkRows(h, f,
		"Ulli, Haupt",
		"Engel, Loritz",
		"Maria, Spannagel",
		"Sigi, Goldschmidt",
		"Engel, Kempf",
		"Manu, Specht",
		"Ulli, Förstner")
}

func Test9006_filtering_fnc_ast_expr(t *testing.T) {
	var (
		ctx, h, s = setup(t)
		m, _, dd  = loadScenario(ctx, s, t, h)
		ff        = loadNoErrMulti(ctx, h, m, dd...)
		f         *report.Frame
	)

	h.a.Len(ff, 1)

	f = ff[0]
	h.a.Len(f.Rows, 7)
	h.a.Equal("first_name<String>, last_name<String>", f.Columns.String())
	checkRows(h, f,
		"Ulli, Haupt",
		"Engel, Loritz",
		"Maria, Spannagel",
		"Sigi, Goldschmidt",
		"Engel, Kempf",
		"Manu, Specht",
		"Ulli, Förstner")
}
