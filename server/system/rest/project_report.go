package rest

import (
	"context"

	"github.com/crusttech/human/server/system/rest/request"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
)

type (
	ProjectReport struct {
		reportSvc projectReportService
	}

	projectReportService interface {
		Report(ctx context.Context, rr *types.ProjectReportRequest) (*types.ProjectReportResult, error)
	}
)

func (ProjectReport) New() *ProjectReport {
	return &ProjectReport{
		reportSvc: service.DefaultProjectReport,
	}
}

func (ctrl *ProjectReport) Report(ctx context.Context, r *request.ProjectReportReport) (interface{}, error) {
	return ctrl.reportSvc.Report(ctx, &types.ProjectReportRequest{
		Resource:      r.Resource,
		ProjectID:     r.ProjectID,
		Dimensions:    r.Dimensions,
		Metrics:       r.Metrics,
		FromTimestamp: r.From,
		ToTimestamp:   r.To,
	})
}
