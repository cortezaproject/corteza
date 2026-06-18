package rest

import (
	"context"

	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/system/reporting"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
)

type (
	Report struct {
		report reportService
		ac     reportAccessController
	}

	reportService interface {
		FindByID(ctx context.Context, ID uint64) (app *types.Report, err error)
		Search(ctx context.Context, filter types.ReportFilter) (aa types.ReportSet, f types.ReportFilter, err error)
		Create(ctx context.Context, new *types.Report) (app *types.Report, err error)
		Update(ctx context.Context, upd *types.Report) (app *types.Report, err error)
		DeleteByID(ctx context.Context, ID uint64) (err error)
		UndeleteByID(ctx context.Context, ID uint64) (err error)
		Run(ctx context.Context, ID uint64, dd reporting.FrameDefinitionSet) (rr []*reporting.Frame, err error)
		Describe(ctx context.Context, src types.ReportDataSourceSet, st types.ReportStepSet, sources ...string) (out []reporting.FrameDescription, err error)
	}

	reportAccessController interface {
		CanGrant(context.Context) bool

		CanReadReport(context.Context, *types.Report) bool
		CanUpdateReport(context.Context, *types.Report) bool
		CanDeleteReport(context.Context, *types.Report) bool
		CanRunReport(context.Context, *types.Report) bool
	}

	reportPayload struct {
		*types.Report

		CanGrant        bool `json:"canGrant"`
		CanReadReport   bool `json:"canReadReport"`
		CanUpdateReport bool `json:"canUpdateReport"`
		CanDeleteReport bool `json:"canDeleteReport"`
		CanRunReport    bool `json:"canRunReport"`
	}

	reportSetPayload struct {
		Filter types.ReportFilter `json:"filter"`
		Set    []*reportPayload   `json:"set"`
	}

	reportFramePayload struct {
		Frames []*reporting.Frame `json:"frames"`
	}
)

func (Report) New() *Report {
	return &Report{
		report: service.DefaultReport,
		ac:     service.DefaultAccessControl,
	}
}

// makeFilter builds the search filter for the generated List controller.
func (ctrl Report) makeFilter(ctx context.Context, r *request.ReportList) (types.ReportFilter, error) {
	var (
		err error
		f   = types.ReportFilter{
			Handle:  r.Handle,
			Query:   r.Query,
			Labels:  r.Labels,
			Deleted: filter.State(r.Deleted),
		}
	)

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return f, err
	}

	f.IncTotal = r.IncTotal

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return f, err
	}

	return f, nil
}

// beforeCreate fills the complex/hook-managed fields onto the resource
// before it is handed to the service. The generated Create controller
// already mapped the plain-value params.
func (ctrl Report) beforeCreate(ctx context.Context, res *types.Report, r *request.ReportCreate) error {
	res.Meta = r.Meta
	res.Scenarios = r.Scenarios
	res.Sources = r.Sources
	res.Blocks = r.Blocks
	res.Labels = r.Labels
	return nil
}

// beforeUpdate fills the complex/hook-managed fields onto the resource
// before it is handed to the service. The generated Update controller
// already mapped the plain-value params (and ID/UpdatedAt).
func (ctrl Report) beforeUpdate(ctx context.Context, res *types.Report, r *request.ReportUpdate) error {
	res.Meta = r.Meta
	res.Scenarios = r.Scenarios
	res.Sources = r.Sources
	res.Blocks = r.Blocks
	res.Labels = r.Labels
	return nil
}

func (ctrl *Report) Describe(ctx context.Context, r *request.ReportDescribe) (interface{}, error) {
	return ctrl.report.Describe(ctx, r.Sources, r.Steps, r.Describe...)
}

func (ctrl *Report) Run(ctx context.Context, r *request.ReportRun) (interface{}, error) {
	rr, err := ctrl.report.Run(ctx, r.ReportID, r.Frames)
	return ctrl.makeReportFramePayload(ctx, rr, err)
}

func (ctrl Report) makePayload(ctx context.Context, m *types.Report, err error) (*reportPayload, error) {
	if err != nil || m == nil {
		return nil, err
	}

	return &reportPayload{
		Report: m,

		CanGrant: ctrl.ac.CanGrant(ctx),

		CanReadReport:   ctrl.ac.CanReadReport(ctx, m),
		CanUpdateReport: ctrl.ac.CanUpdateReport(ctx, m),
		CanDeleteReport: ctrl.ac.CanDeleteReport(ctx, m),
		CanRunReport:    ctrl.ac.CanRunReport(ctx, m),
	}, nil
}

func (ctrl Report) makeFilterPayload(ctx context.Context, nn types.ReportSet, f types.ReportFilter, err error) (*reportSetPayload, error) {
	if err != nil {
		return nil, err
	}

	msp := &reportSetPayload{Filter: f, Set: make([]*reportPayload, len(nn))}

	for i := range nn {
		msp.Set[i], _ = ctrl.makePayload(ctx, nn[i], nil)
	}

	return msp, nil
}

func (ctrl Report) makeReportFramePayload(ctx context.Context, ff []*reporting.Frame, err error) (*reportFramePayload, error) {
	if err != nil || len(ff) == 0 {
		return nil, err
	}

	return &reportFramePayload{
		Frames: ff,
	}, nil
}
