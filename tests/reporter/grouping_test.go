package reporter

import (
	"testing"

	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/report"
)

func TestReporterGrouping(t *testing.T) {
	ctx, h, s, rp := prepare(t, "report_grouping_base")
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

	t.Run("basic grouping", func(t *testing.T) {
		rr, err := model.Load(ctx, fd)
		h.a.NoError(err)
		h.a.Len(rr, 1)
		r := rr[0]
		h.a.Equal(6, r.Size())

		// columns
		h.a.Equal("by_name<String>, count<Number>, total<Number>", r.Columns.String())

		req := []string{
			"Engel, 3, 179",
			"Manu, 1, 61",
			"Maria, 3, 183",
			"Sascha, 1, 38",
			"Sigi, 1, 67",
			"Ulli, 3, 122",
		}

		// rows
		r.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Equal(req[i], r.String())
			return nil
		})
	})

	t.Run("basic filtering", func(t *testing.T) {
		fd.Rows = &report.RowDefinition{
			And: []*report.RowDefinition{
				{
					Cells: map[string]*report.CellDefinition{
						"total": {Op: "gt", Value: "50"},
						"count": {Op: "lt", Value: "2"},
					},
				},
			},
		}

		rr, err := model.Load(ctx, fd)
		h.a.NoError(err)
		h.a.Len(rr, 1)
		r := rr[0]
		h.a.Equal(2, r.Size())

		// columns
		h.a.Equal("by_name<String>, count<Number>, total<Number>", r.Columns.String())

		req := []string{
			"Manu, 1, 61",
			"Sigi, 1, 67",
		}

		// rows
		r.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Equal(req[i], r.String())
			return nil
		})
	})

	t.Run("basic sorting", func(t *testing.T) {
		fd.Rows = nil
		fd.Sorting = filter.SortExprSet{
			{Column: "count", Descending: true},
			{Column: "by_name", Descending: false},
		}

		rr, err := model.Load(ctx, fd)
		h.a.NoError(err)
		h.a.Len(rr, 1)
		r := rr[0]
		h.a.Equal(6, r.Size())

		// columns
		h.a.Equal("by_name<String>, count<Number>, total<Number>", r.Columns.String())

		req := []string{
			"Engel, 3, 179",
			"Maria, 3, 183",
			"Ulli, 3, 122",
			"Manu, 1, 61",
			"Sascha, 1, 38",
			"Sigi, 1, 67",
		}

		// rows
		r.WalkRows(func(i int, r report.FrameRow) error {
			h.a.Equal(req[i], r.String())
			return nil
		})
	})
}
