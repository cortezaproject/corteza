package rest

import (
	"context"
	"github.com/cortezaproject/corteza-server/system/rest/request"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/pkg/errors"
)

var _ = errors.Wrap

type (
	Stats struct {
		svc statsService
	}

	statsService interface {
		Metrics(context.Context) (*service.StatisticsMetricsPayload, error)
	}
)

func (Stats) New() *Stats {
	return &Stats{
		svc: service.DefaultStatistics,
	}
}

func (ctrl *Stats) List(ctx context.Context, r *request.StatsList) (interface{}, error) {
	return ctrl.svc.Metrics(ctx)
}
