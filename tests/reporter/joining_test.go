package reporter

import (
	"testing"

	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/report"
)

func TestReporterJoining(t *testing.T) {
	ctx, h, s, rp := prepare(t, "report_joining_base")
	providers := map[string]report.DatasourceProvider{
		"composeRecords": service.DefaultRecord,
	}
	_ = s

	ss := rp.Sources.ModelSteps()
	model, err := report.Model(ctx, providers, ss...)
	h.a.NoError(err)
	err = model.Run(ctx)
	h.a.NoError(err)

	t.Run("basic loading", func(t *testing.T) {
		rr, err := model.Load(ctx, rp.Frames...)
		h.a.NoError(err)
		h.a.Len(rr, 7)
		local := rr[0]
		ix := indexJoinedResult(rr)
		_ = ix

		// local
		h.a.Equal(12, local.Size())
		h.a.Equal("id<Record>, join_key<String>, first_name<String>, last_name<String>", local.Columns.String())
		req := []string{
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
			", Ulli_Förstner, Ulli, Förstner",
		}
		local.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Contains(r.String(), req[i])
			return nil
		})

		// Maria_Königsmann
		f := ix["Maria_Königsmann"]
		h.a.NotNil(f)
		h.a.Equal("id<Record>, usr<String>, name<String>, type<Select>, cost<Number>, time_spent<Number>", f.Columns.String())
		req = []string{
			", Maria_Königsmann, u1 j1, a, 10, 2",
			", Maria_Königsmann, u1 j5, a, 4, 4",
			", Maria_Königsmann, u1 j2, b, 20, 5",
			", Maria_Königsmann, u1 j6, b, 25, 25",
			", Maria_Königsmann, u1 j7, b, 9, 91",
			", Maria_Königsmann, u1 j4, c, 3, 0",
			", Maria_Königsmann, u1 j3, d, 11, 1",
		}
		f.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Contains(r.String(), req[i])
			return nil
		})

		// Engel_Loritz
		f = ix["Engel_Loritz"]
		h.a.NotNil(f)
		h.a.Equal("id<Record>, usr<String>, name<String>, type<Select>, cost<Number>, time_spent<Number>", f.Columns.String())
		req = []string{
			", Engel_Loritz, u3 j1, a, 10, 1",
			", Engel_Loritz, u3 j2, a, 0, 0",
			", Engel_Loritz, u3 j3, a, 19, 99",
		}
		f.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Contains(r.String(), req[i])
			return nil
		})

		// Sigi_Goldschmidt
		f = ix["Sigi_Goldschmidt"]
		h.a.NotNil(f)
		h.a.Equal("id<Record>, usr<String>, name<String>, type<Select>, cost<Number>, time_spent<Number>", f.Columns.String())
		req = []string{
			", Sigi_Goldschmidt, u7 j2, a, 10, 21",
			", Sigi_Goldschmidt, u7 j3, b, 10, 99",
			", Sigi_Goldschmidt, u7 j1, d, 10, 29",
		}
		f.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Contains(r.String(), req[i])
			return nil
		})

		// Engel_Kiefer
		f = ix["Engel_Kiefer"]
		h.a.NotNil(f)
		h.a.Equal("id<Record>, usr<String>, name<String>, type<Select>, cost<Number>, time_spent<Number>", f.Columns.String())
		req = []string{
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
			", Engel_Kiefer, u12 j7, d, 14, 71",
		}
		f.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Contains(r.String(), req[i])
			return nil
		})

		// Manu_Specht
		f = ix["Manu_Specht"]
		h.a.NotNil(f)
		h.a.Equal("id<Record>, usr<String>, name<String>, type<Select>, cost<Number>, time_spent<Number>", f.Columns.String())
		req = []string{
			", Manu_Specht, u10 j3, b, 53, 12",
			", Manu_Specht, u10 j4, b, 60, 22",
			", Manu_Specht, u10 j2, c, 83, 70",
			", Manu_Specht, u10 j1, d, 45, 56",
		}
		f.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Contains(r.String(), req[i])
			return nil
		})

		// Ulli_Böhler
		f = ix["Ulli_Böhler"]
		h.a.NotNil(f)
		h.a.Equal("id<Record>, usr<String>, name<String>, type<Select>, cost<Number>, time_spent<Number>", f.Columns.String())
		req = []string{
			", Ulli_Böhler, u5 j1, a, 1, 2",
		}
		f.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Contains(r.String(), req[i])
			return nil
		})
	})

	t.Run("basic filtering", func(t *testing.T) {
		// users
		rp.Frames[0].Rows = &report.RowDefinition{
			Or: []*report.RowDefinition{
				{
					Cells: map[string]*report.CellDefinition{
						"last_name": {Op: "eq", Value: "'Königsmann'"},
					},
				},
				{
					Cells: map[string]*report.CellDefinition{
						"last_name": {Op: "eq", Value: "'Kiefer'"},
					},
				},
			},
		}
		// jobs
		rp.Frames[1].Rows = &report.RowDefinition{
			And: []*report.RowDefinition{
				{
					Cells: map[string]*report.CellDefinition{
						"type": {Op: "eq", Value: "'a'"},
					},
				},
			},
		}

		rr, err := model.Load(ctx, rp.Frames...)
		h.a.NoError(err)
		h.a.Len(rr, 3)
		local := rr[0]
		ix := indexJoinedResult(rr)
		_ = ix

		// local
		h.a.Equal(2, local.Size())
		h.a.Equal("id<Record>, join_key<String>, first_name<String>, last_name<String>", local.Columns.String())
		req := []string{
			", Engel_Kiefer, Engel, Kiefer",
			", Maria_Königsmann, Maria, Königsmann",
		}
		local.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Contains(r.String(), req[i])
			return nil
		})

		// Engel_Kiefer
		f := ix["Engel_Kiefer"]
		h.a.NotNil(f)
		h.a.Equal(4, f.Size())
		h.a.Equal("id<Record>, usr<String>, name<String>, type<Select>, cost<Number>, time_spent<Number>", f.Columns.String())
		req = []string{
			", Engel_Kiefer, u12 j1, a, 42, 69",
			", Engel_Kiefer, u12 j4, a, 35, 26",
			", Engel_Kiefer, u12 j5, a, 34, 29",
			", Engel_Kiefer, u12 j9, a, 71, 90",
		}
		f.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Contains(r.String(), req[i])
			return nil
		})

		// Engel_Loritz
		f = ix["Maria_Königsmann"]
		h.a.NotNil(f)
		h.a.Equal(2, f.Size())
		h.a.Equal("id<Record>, usr<String>, name<String>, type<Select>, cost<Number>, time_spent<Number>", f.Columns.String())
		req = []string{
			", Maria_Königsmann, u1 j1, a, 10, 2",
			", Maria_Königsmann, u1 j5, a, 4, 4",
		}
		f.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Contains(r.String(), req[i])
			return nil
		})
	})

	t.Run("basic sorting; local datasource", func(t *testing.T) {
		// users
		rp.Frames[0].Rows = nil
		rp.Frames[0].Sorting = filter.SortExprSet{
			{Column: "first_name", Descending: false},
			{Column: "last_name", Descending: true},
		}
		// jobs
		rp.Frames[1].Rows = nil

		rr, err := model.Load(ctx, rp.Frames...)
		h.a.NoError(err)
		h.a.Len(rr, 7)
		local := rr[0]
		ix := indexJoinedResult(rr)
		_ = ix

		// local
		h.a.Equal(12, local.Size())
		h.a.Equal("id<Record>, join_key<String>, first_name<String>, last_name<String>", local.Columns.String())
		req := []string{
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
			", Ulli, Böhler",
		}
		local.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Contains(r.String(), req[i])
			return nil
		})

		// Maria_Königsmann
		f := ix["Maria_Königsmann"]
		h.a.NotNil(f)
		h.a.Equal("id<Record>, usr<String>, name<String>, type<Select>, cost<Number>, time_spent<Number>", f.Columns.String())
		req = []string{
			", Maria_Königsmann, u1 j1, a, 10, 2",
			", Maria_Königsmann, u1 j5, a, 4, 4",
			", Maria_Königsmann, u1 j2, b, 20, 5",
			", Maria_Königsmann, u1 j6, b, 25, 25",
			", Maria_Königsmann, u1 j7, b, 9, 91",
			", Maria_Königsmann, u1 j4, c, 3, 0",
			", Maria_Königsmann, u1 j3, d, 11, 1",
		}
		f.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Contains(r.String(), req[i])
			return nil
		})

		// Engel_Loritz
		f = ix["Engel_Loritz"]
		h.a.NotNil(f)
		h.a.Equal("id<Record>, usr<String>, name<String>, type<Select>, cost<Number>, time_spent<Number>", f.Columns.String())
		req = []string{
			", Engel_Loritz, u3 j1, a, 10, 1",
			", Engel_Loritz, u3 j2, a, 0, 0",
			", Engel_Loritz, u3 j3, a, 19, 99",
		}
		f.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Contains(r.String(), req[i])
			return nil
		})

		// Sigi_Goldschmidt
		f = ix["Sigi_Goldschmidt"]
		h.a.NotNil(f)
		h.a.Equal("id<Record>, usr<String>, name<String>, type<Select>, cost<Number>, time_spent<Number>", f.Columns.String())
		req = []string{
			", Sigi_Goldschmidt, u7 j2, a, 10, 21",
			", Sigi_Goldschmidt, u7 j3, b, 10, 99",
			", Sigi_Goldschmidt, u7 j1, d, 10, 29",
		}
		f.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Contains(r.String(), req[i])
			return nil
		})

		// Engel_Kiefer
		f = ix["Engel_Kiefer"]
		h.a.NotNil(f)
		h.a.Equal("id<Record>, usr<String>, name<String>, type<Select>, cost<Number>, time_spent<Number>", f.Columns.String())
		req = []string{
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
			", Engel_Kiefer, u12 j7, d, 14, 71",
		}
		f.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Contains(r.String(), req[i])
			return nil
		})

		// Manu_Specht
		f = ix["Manu_Specht"]
		h.a.NotNil(f)
		h.a.Equal("id<Record>, usr<String>, name<String>, type<Select>, cost<Number>, time_spent<Number>", f.Columns.String())
		req = []string{
			", Manu_Specht, u10 j3, b, 53, 12",
			", Manu_Specht, u10 j4, b, 60, 22",
			", Manu_Specht, u10 j2, c, 83, 70",
			", Manu_Specht, u10 j1, d, 45, 56",
		}
		f.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Contains(r.String(), req[i])
			return nil
		})

		// Ulli_Böhler
		f = ix["Ulli_Böhler"]
		h.a.NotNil(f)
		h.a.Equal("id<Record>, usr<String>, name<String>, type<Select>, cost<Number>, time_spent<Number>", f.Columns.String())
		req = []string{
			", Ulli_Böhler, u5 j1, a, 1, 2",
		}
		f.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Contains(r.String(), req[i])
			return nil
		})
	})

	t.Run("basic sorting; foreign datasource", func(t *testing.T) {
		// @todo !!! for now local DS sorting based on the foreign DS includes
		//       only the rows that have a match.
		//       This will do for now, but should be improved in the future.

		// users
		rp.Frames[0].Rows = nil
		rp.Frames[0].Sorting = filter.SortExprSet{
			{Column: "jobs.type", Descending: false},
			{Column: "first_name", Descending: false},
			{Column: "last_name", Descending: true},
		}
		// jobs
		rp.Frames[1].Rows = nil
		rp.Frames[1].Sorting = nil

		rr, err := model.Load(ctx, rp.Frames...)
		h.a.NoError(err)
		h.a.Len(rr, 7)
		local := rr[0]
		ix := indexJoinedResult(rr)
		_ = ix

		// local
		h.a.Equal(6, local.Size())
		h.a.Equal("id<Record>, join_key<String>, first_name<String>, last_name<String>", local.Columns.String())
		req := []string{
			", Engel, Loritz",
			", Engel, Kiefer",
			", Maria, Königsmann",
			", Sigi, Goldschmidt",
			", Ulli, Böhler",
			", Manu, Specht",
		}
		local.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Contains(r.String(), req[i])
			return nil
		})

		// Maria_Königsmann
		f := ix["Maria_Königsmann"]
		h.a.NotNil(f)
		h.a.Equal("id<Record>, usr<String>, name<String>, type<Select>, cost<Number>, time_spent<Number>", f.Columns.String())
		req = []string{
			", Maria_Königsmann, u1 j1, a, 10, 2",
			", Maria_Königsmann, u1 j5, a, 4, 4",
			", Maria_Königsmann, u1 j2, b, 20, 5",
			", Maria_Königsmann, u1 j6, b, 25, 25",
			", Maria_Königsmann, u1 j7, b, 9, 91",
			", Maria_Königsmann, u1 j4, c, 3, 0",
			", Maria_Königsmann, u1 j3, d, 11, 1",
		}
		f.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Contains(r.String(), req[i])
			return nil
		})

		// Engel_Loritz
		f = ix["Engel_Loritz"]
		h.a.NotNil(f)
		h.a.Equal("id<Record>, usr<String>, name<String>, type<Select>, cost<Number>, time_spent<Number>", f.Columns.String())
		req = []string{
			", Engel_Loritz, u3 j1, a, 10, 1",
			", Engel_Loritz, u3 j2, a, 0, 0",
			", Engel_Loritz, u3 j3, a, 19, 99",
		}
		f.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Contains(r.String(), req[i])
			return nil
		})

		// Sigi_Goldschmidt
		f = ix["Sigi_Goldschmidt"]
		h.a.NotNil(f)
		h.a.Equal("id<Record>, usr<String>, name<String>, type<Select>, cost<Number>, time_spent<Number>", f.Columns.String())
		req = []string{
			", Sigi_Goldschmidt, u7 j2, a, 10, 21",
			", Sigi_Goldschmidt, u7 j3, b, 10, 99",
			", Sigi_Goldschmidt, u7 j1, d, 10, 29",
		}
		f.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Contains(r.String(), req[i])
			return nil
		})

		// Engel_Kiefer
		f = ix["Engel_Kiefer"]
		h.a.NotNil(f)
		h.a.Equal("id<Record>, usr<String>, name<String>, type<Select>, cost<Number>, time_spent<Number>", f.Columns.String())
		req = []string{
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
			", Engel_Kiefer, u12 j7, d, 14, 71",
		}
		f.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Contains(r.String(), req[i])
			return nil
		})

		// Manu_Specht
		f = ix["Manu_Specht"]
		h.a.NotNil(f)
		h.a.Equal("id<Record>, usr<String>, name<String>, type<Select>, cost<Number>, time_spent<Number>", f.Columns.String())
		req = []string{
			", Manu_Specht, u10 j3, b, 53, 12",
			", Manu_Specht, u10 j4, b, 60, 22",
			", Manu_Specht, u10 j2, c, 83, 70",
			", Manu_Specht, u10 j1, d, 45, 56",
		}
		f.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Contains(r.String(), req[i])
			return nil
		})

		// Ulli_Böhler
		f = ix["Ulli_Böhler"]
		h.a.NotNil(f)
		h.a.Equal("id<Record>, usr<String>, name<String>, type<Select>, cost<Number>, time_spent<Number>", f.Columns.String())
		req = []string{
			", Ulli_Böhler, u5 j1, a, 1, 2",
		}
		f.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Contains(r.String(), req[i])
			return nil
		})
	})

	t.Run("paging", func(t *testing.T) {
		rp.Frames[0].Paging = &filter.Paging{
			Limit: 5,
		}
		rp.Frames[0].Sorting = filter.SortExprSet{
			&filter.SortExpr{Column: "join_key", Descending: false},
		}

		rp.Frames[1].Paging = &filter.Paging{
			Limit: 2,
		}
		rp.Frames[1].Sorting = filter.SortExprSet{
			&filter.SortExpr{Column: "name", Descending: false},
		}

		// ^ going up ^

		// // // PAGE 1
		rr, err := model.Load(ctx, rp.Frames...)
		h.a.NoError(err)
		h.a.Len(rr, 4)
		local := rr[0]
		ix := indexJoinedResult(rr)
		_ = ix

		// local
		h.a.Equal(5, local.Size())
		h.a.NotNil(local.Paging)
		h.a.NotNil(local.Paging.NextPage)
		h.a.Nil(local.Paging.PrevPage)
		req := []string{
			", Engel_Kempf, Engel, Kempf",
			", Engel_Kiefer, Engel, Kiefer",
			", Engel_Loritz, Engel, Loritz",
			", Manu_Specht, Manu, Specht",
			", Maria_Krüger, Maria, Krüger",
		}
		local.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Contains(r.String(), req[i])
			return nil
		})

		// Manu_Specht
		f := ix["Manu_Specht"]
		h.a.NotNil(f)
		h.a.NotNil(f.Paging)
		h.a.NotNil(f.Paging.NextPage)
		h.a.Nil(f.Paging.PrevPage)
		req = []string{
			", Manu_Specht, u10 j1, d, 45, 56",
			", Manu_Specht, u10 j2, c, 83, 70",
		}
		f.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Contains(r.String(), req[i])
			return nil
		})

		// Engel_Kiefer
		f = ix["Engel_Kiefer"]
		h.a.NotNil(f)
		h.a.NotNil(f.Paging)
		h.a.NotNil(f.Paging.NextPage)
		h.a.Nil(f.Paging.PrevPage)
		req = []string{
			", Engel_Kiefer, u12 j1, a, 42, 69",
			", Engel_Kiefer, u12 j10, c, 79, 25",
		}
		f.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Contains(r.String(), req[i])
			return nil
		})

		// Engel_Loritz
		f = ix["Engel_Loritz"]
		h.a.NotNil(f)
		h.a.NotNil(f.Paging)
		h.a.NotNil(f.Paging.NextPage)
		h.a.Nil(f.Paging.PrevPage)
		req = []string{
			", Engel_Loritz, u3 j1, a, 10, 1",
			", Engel_Loritz, u3 j2, a, 0, 0",
		}
		f.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Contains(r.String(), req[i])
			return nil
		})

		// // // PAGE 2
		rp.Frames[0].Paging.PageCursor = local.Paging.NextPage
		rr, err = model.Load(ctx, rp.Frames...)
		h.a.NoError(err)

		h.a.Len(rr, 4)
		local = rr[0]
		ix = indexJoinedResult(rr)
		_ = ix

		// local
		h.a.Equal(5, local.Size())
		h.a.NotNil(local.Paging)
		h.a.NotNil(local.Paging.NextPage)
		h.a.NotNil(local.Paging.PrevPage)
		req = []string{
			", Maria_Königsmann, Maria, Königsmann",
			", Maria_Spannagel, Maria, Spannagel",
			", Sascha_Jans, Sascha, Jans",
			", Sigi_Goldschmidt, Sigi, Goldschmidt",
			", Ulli_Böhler, Ulli, Böhler",
			", Ulli_Förstner, Ulli, Förstner",
		}
		local.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Contains(r.String(), req[i])
			return nil
		})

		// Maria_Königsmann
		f = ix["Maria_Königsmann"]
		h.a.NotNil(f)
		h.a.NotNil(f.Paging)
		h.a.NotNil(f.Paging.NextPage)
		h.a.Nil(f.Paging.PrevPage)
		req = []string{
			", Maria_Königsmann, u1 j1, a, 10, 2",
			", Maria_Königsmann, u1 j2, b, 20, 5",
		}
		f.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Contains(r.String(), req[i])
			return nil
		})

		// Engel_Kiefer
		f = ix["Ulli_Böhler"]
		h.a.NotNil(f)
		h.a.Nil(f.Paging)
		req = []string{
			", Ulli_Böhler, u5 j1, a, 1, 2",
		}
		f.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Contains(r.String(), req[i])
			return nil
		})

		// Sigi_Goldschmidt
		f = ix["Sigi_Goldschmidt"]
		h.a.NotNil(f)
		h.a.NotNil(f.Paging)
		h.a.NotNil(f.Paging.NextPage)
		h.a.Nil(f.Paging.PrevPage)
		req = []string{
			", Sigi_Goldschmidt, u7 j1, d, 10, 29",
			", Sigi_Goldschmidt, u7 j2, a, 10, 21",
		}
		f.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Contains(r.String(), req[i])
			return nil
		})

		// // // PAGE 3
		rp.Frames[0].Paging.PageCursor = local.Paging.NextPage
		rr, err = model.Load(ctx, rp.Frames...)
		h.a.NoError(err)

		h.a.Len(rr, 1)
		local = rr[0]
		ix = indexJoinedResult(rr)
		_ = ix

		// local
		h.a.Equal(2, local.Size())
		h.a.NotNil(local.Paging)
		h.a.Nil(local.Paging.NextPage)
		h.a.NotNil(local.Paging.PrevPage)
		req = []string{
			", Ulli_Förstner, Ulli, Förstner",
			", Ulli_Haupt, Ulli, Haupt",
		}
		local.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Contains(r.String(), req[i])
			return nil
		})

		// v going down v

		// // // PAGE 2
		rp.Frames[0].Paging.PageCursor = local.Paging.PrevPage
		rr, err = model.Load(ctx, rp.Frames...)
		h.a.NoError(err)

		h.a.Len(rr, 4)
		local = rr[0]
		ix = indexJoinedResult(rr)
		_ = ix

		// local
		h.a.Equal(5, local.Size())
		h.a.NotNil(local.Paging)
		h.a.NotNil(local.Paging.NextPage)
		h.a.NotNil(local.Paging.PrevPage)
		req = []string{
			", Maria_Königsmann, Maria, Königsmann",
			", Maria_Spannagel, Maria, Spannagel",
			", Sascha_Jans, Sascha, Jans",
			", Sigi_Goldschmidt, Sigi, Goldschmidt",
			", Ulli_Böhler, Ulli, Böhler",
			", Ulli_Förstner, Ulli, Förstner",
		}
		local.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Contains(r.String(), req[i])
			return nil
		})

		// Maria_Königsmann
		f = ix["Maria_Königsmann"]
		h.a.NotNil(f)
		h.a.NotNil(f.Paging)
		h.a.NotNil(f.Paging.NextPage)
		h.a.Nil(f.Paging.PrevPage)
		req = []string{
			", Maria_Königsmann, u1 j1, a, 10, 2",
			", Maria_Königsmann, u1 j2, b, 20, 5",
		}
		f.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Contains(r.String(), req[i])
			return nil
		})

		// Engel_Kiefer
		f = ix["Ulli_Böhler"]
		h.a.NotNil(f)
		h.a.Nil(f.Paging)
		req = []string{
			", Ulli_Böhler, u5 j1, a, 1, 2",
		}
		f.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Contains(r.String(), req[i])
			return nil
		})

		// Sigi_Goldschmidt
		f = ix["Sigi_Goldschmidt"]
		h.a.NotNil(f)
		h.a.NotNil(f.Paging)
		h.a.NotNil(f.Paging.NextPage)
		h.a.Nil(f.Paging.PrevPage)
		req = []string{
			", Sigi_Goldschmidt, u7 j1, d, 10, 29",
			", Sigi_Goldschmidt, u7 j2, a, 10, 21",
		}
		f.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Contains(r.String(), req[i])
			return nil
		})

		// // // PAGE 1
		rp.Frames[0].Paging.PageCursor = local.Paging.PrevPage
		rr, err = model.Load(ctx, rp.Frames...)
		h.a.NoError(err)
		h.a.Len(rr, 4)
		local = rr[0]
		ix = indexJoinedResult(rr)
		_ = ix

		// local
		h.a.Equal(5, local.Size())
		h.a.NotNil(local.Paging)
		h.a.NotNil(local.Paging.NextPage)
		h.a.Nil(local.Paging.PrevPage)
		req = []string{
			", Engel_Kempf, Engel, Kempf",
			", Engel_Kiefer, Engel, Kiefer",
			", Engel_Loritz, Engel, Loritz",
			", Manu_Specht, Manu, Specht",
			", Maria_Krüger, Maria, Krüger",
		}
		local.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Contains(r.String(), req[i])
			return nil
		})

		// Manu_Specht
		f = ix["Manu_Specht"]
		h.a.NotNil(f)
		h.a.NotNil(f.Paging)
		h.a.NotNil(f.Paging.NextPage)
		h.a.Nil(f.Paging.PrevPage)
		req = []string{
			", Manu_Specht, u10 j1, d, 45, 56",
			", Manu_Specht, u10 j2, c, 83, 70",
		}
		f.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Contains(r.String(), req[i])
			return nil
		})

		// Engel_Kiefer
		f = ix["Engel_Kiefer"]
		h.a.NotNil(f)
		h.a.NotNil(f.Paging)
		h.a.NotNil(f.Paging.NextPage)
		h.a.Nil(f.Paging.PrevPage)
		req = []string{
			", Engel_Kiefer, u12 j1, a, 42, 69",
			", Engel_Kiefer, u12 j10, c, 79, 25",
		}
		f.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Contains(r.String(), req[i])
			return nil
		})

		// Engel_Loritz
		f = ix["Engel_Loritz"]
		h.a.NotNil(f)
		h.a.NotNil(f.Paging)
		h.a.NotNil(f.Paging.NextPage)
		h.a.Nil(f.Paging.PrevPage)
		req = []string{
			", Engel_Loritz, u3 j1, a, 10, 1",
			", Engel_Loritz, u3 j2, a, 0, 0",
		}
		f.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Contains(r.String(), req[i])
			return nil
		})
	})
}

func indexJoinedResult(ff []*report.Frame) map[string]*report.Frame {
	out := make(map[string]*report.Frame)
	// the first one is the local ds
	for _, f := range ff[1:] {
		out[f.RefValue] = f
	}

	return out
}
