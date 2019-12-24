package rest

import (
	"context"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/system/rest/request"
	"github.com/cortezaproject/corteza-server/system/service"
)

var _ = errors.Wrap

type (
	Stats struct {
		svc statsService
	}

	statsService interface {
		Metrics(context.Context) (interface{}, error)
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
