package rest

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/report"
	"github.com/cortezaproject/corteza-server/system/rest/request"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/pkg/errors"
)

var _ = errors.Wrap

type (
	Report struct {
		report reportService
		ac     reportAccessController
	}

	reportService interface {
		LookupByID(ctx context.Context, ID uint64) (app *types.Report, err error)
		Search(ctx context.Context, filter types.ReportFilter) (aa types.ReportSet, f types.ReportFilter, err error)
		Create(ctx context.Context, new *types.Report) (app *types.Report, err error)
		Update(ctx context.Context, upd *types.Report) (app *types.Report, err error)
		Delete(ctx context.Context, ID uint64) (err error)
		Undelete(ctx context.Context, ID uint64) (err error)
		Run(ctx context.Context, ID uint64, dd report.FrameDefinitionSet) (rr []*report.Frame, err error)
		DescribeFresh(ctx context.Context, src types.ReportDataSourceSet, st report.StepDefinitionSet, sources ...string) (out report.FrameDescriptionSet, err error)
	}

	reportAccessController interface {
		CanGrant(context.Context) bool

		CanUpdateReport(context.Context, *types.Report) bool
		CanDeleteReport(context.Context, *types.Report) bool
		CanRunReport(context.Context, *types.Report) bool
	}

	reportPayload struct {
		*types.Report

		CanGrant        bool `json:"canGrant"`
		CanUpdateReport bool `json:"canUpdateReport"`
		CanDeleteReport bool `json:"canDeleteReport"`
		CanRunReport    bool `json:"canRunReport"`
	}

	reportSetPayload struct {
		Filter types.ReportFilter `json:"filter"`
		Set    []*reportPayload   `json:"set"`
	}

	reportFramePayload struct {
		Frames []*report.Frame `json:"frames"`
	}
)

func (Report) New() *Report {
	return &Report{
		report: service.DefaultReport,
		ac:     service.DefaultAccessControl,
	}
}

func (ctrl *Report) List(ctx context.Context, r *request.ReportList) (interface{}, error) {
	var (
		err error
		f   = types.ReportFilter{
			Handle:  r.Handle,
			Labels:  r.Labels,
			Deleted: filter.State(r.Deleted),
		}
	)

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return nil, err
	}

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return nil, err
	}

	set, filter, err := ctrl.report.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *Report) Create(ctx context.Context, r *request.ReportCreate) (interface{}, error) {
	var (
		err error
		app = &types.Report{
			Handle:  r.Handle,
			Meta:    r.Meta,
			Sources: r.Sources,
			Blocks:  r.Blocks,
			Labels:  r.Labels,
		}
	)

	app, err = ctrl.report.Create(ctx, app)
	return ctrl.makePayload(ctx, app, err)
}

func (ctrl *Report) Update(ctx context.Context, r *request.ReportUpdate) (interface{}, error) {
	var (
		err error
		app = &types.Report{
			ID:      r.ReportID,
			Handle:  r.Handle,
			Meta:    r.Meta,
			Sources: r.Sources,
			Blocks:  r.Blocks,
			Labels:  r.Labels,
		}
	)

	app, err = ctrl.report.Update(ctx, app)
	return ctrl.makePayload(ctx, app, err)
}

func (ctrl *Report) Read(ctx context.Context, r *request.ReportRead) (interface{}, error) {
	app, err := ctrl.report.LookupByID(ctx, r.ReportID)
	return ctrl.makePayload(ctx, app, err)
}

func (ctrl *Report) Delete(ctx context.Context, r *request.ReportDelete) (interface{}, error) {
	return api.OK(), ctrl.report.Delete(ctx, r.ReportID)
}

func (ctrl *Report) Undelete(ctx context.Context, r *request.ReportUndelete) (interface{}, error) {
	return api.OK(), ctrl.report.Undelete(ctx, r.ReportID)
}

func (ctrl *Report) Describe(ctx context.Context, r *request.ReportDescribe) (interface{}, error) {
	return ctrl.report.DescribeFresh(ctx, r.Sources, r.Steps, r.Describe...)
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

		CanUpdateReport: ctrl.ac.CanUpdateReport(ctx, m),
		CanDeleteReport: ctrl.ac.CanDeleteReport(ctx, m),
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

func (ctrl Report) makeReportFramePayload(ctx context.Context, ff []*report.Frame, err error) (*reportFramePayload, error) {
	if err != nil || len(ff) == 0 {
		return nil, err
	}

	return &reportFramePayload{
		Frames: ff,
	}, nil
}
