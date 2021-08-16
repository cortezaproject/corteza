package reporter

import (
	"context"
	"path"
	"testing"

	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/report"
	"github.com/cortezaproject/corteza-server/store"
)

func TestReporterLoading(t *testing.T) {
	ctx, h, s, rp := prepare(t, "report_loading_base")
	providers := map[string]report.DatasourceProvider{
		"composeRecords": service.DefaultRecord,
	}
	_ = s

	ss := rp.Sources.ModelSteps()
	model, err := report.Model(ctx, providers, ss...)
	h.a.NoError(err)
	err = model.Run(ctx)
	h.a.NoError(err)
	fd := rp.Frames[0]

	t.Run("basic loading", func(t *testing.T) {
		rr, err := model.Load(ctx, fd)
		h.a.NoError(err)
		h.a.Len(rr, 1)
		r := rr[0]
		h.a.Equal(12, r.Size())

		// columns
		h.a.Equal("id<Record>, first_name<String>, last_name<String>, number_of_numbers<Number>", r.Columns.String())

		// rows
		r.WalkRows(func(i int, r report.FrameRow) error {
			for _, c := range r {
				h.a.NotNil(c)
			}

			return nil
		})
	})

	t.Run("basic filtering", func(t *testing.T) {
		fd.Rows = &report.RowDefinition{
			Or: []*report.RowDefinition{
				{
					Cells: map[string]*report.CellDefinition{
						"first_name": {Op: "eq", Value: "'Maria'"},
					},
				},

				{
					// these cells are connected with an OR, because the parent is OR
					Cells: map[string]*report.CellDefinition{
						"first_name": {Op: "eq", Value: "'Ulli'"},
						"last_name":  {Op: "eq", Value: "'Specht'"},
					},
				},
			},
		}

		rr, err := model.Load(ctx, fd)
		h.a.NoError(err)
		h.a.Len(rr, 1)
		r := rr[0]
		// 3xMaria + 3xUlli + 1xSpecht
		h.a.Equal(7, r.Size())

		// columns
		h.a.Equal("id<Record>, first_name<String>, last_name<String>, number_of_numbers<Number>", r.Columns.String())

		// rows
		r.WalkRows(func(i int, r report.FrameRow) error {
			for _, c := range r {
				h.a.NotNil(c)
			}

			return nil
		})
	})

	t.Run("basic sorting", func(t *testing.T) {
		fd.Rows = nil
		fd.Sorting = filter.SortExprSet{
			{Column: "first_name", Descending: false},
			{Column: "last_name", Descending: true},
			{Column: "id", Descending: false},
		}

		rr, err := model.Load(ctx, fd)
		h.a.NoError(err)
		h.a.Len(rr, 1)
		r := rr[0]
		h.a.Equal(12, r.Size())

		// columns
		h.a.Equal("id<Record>, first_name<String>, last_name<String>, number_of_numbers<Number>", r.Columns.String())

		// omit the ID's because they are generated
		req := []string{
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
			", Ulli, Böhler, 14",
		}

		// rows
		r.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Contains(r.String(), req[i])
			return nil
		})
	})
}

func prepare(t *testing.T, report string) (context.Context, helper, store.Storer, *auxReport) {
	h := newHelper(t)
	s := service.DefaultStore
	ctx := auth.SetSuperUserContext(context.Background())

	err := bmReporterPrepareDM(ctx, h, s, "integration")
	h.a.NoError(err)

	rp, err := bmReporterParseReport(path.Join("testdata", "integration", report+".json"))
	h.a.NoError(err)

	return ctx, h, s, rp
}
