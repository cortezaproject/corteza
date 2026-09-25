package rest

import (
	"context"

	"github.com/crusttech/human/server/system/rest/request"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
)

type (
	Stats struct {
		svc statsService
	}

	statsService interface {
		Metrics(context.Context, service.StatisticsRequest) (*types.SystemStats, error)
		Detail(context.Context, string, service.StatisticsRequest) (*types.SystemStatsDetail, error)
		Events(context.Context, string, service.StatisticsRequest) (*types.SystemStatsEvents, error)
	}
)

func (Stats) New() *Stats {
	return &Stats{
		svc: service.DefaultStatistics,
	}
}

func (ctrl *Stats) List(ctx context.Context, r *request.StatsList) (interface{}, error) {
	return ctrl.svc.Metrics(ctx, service.StatisticsRequest{
		From:   r.From,
		To:     r.To,
		Bucket: r.Bucket,
	})
}

func (ctrl *Stats) Detail(ctx context.Context, r *request.StatsDetail) (interface{}, error) {
	return ctrl.svc.Detail(ctx, r.Resource, service.StatisticsRequest{
		From:   r.From,
		To:     r.To,
		Bucket: r.Bucket,
	})
}

func (ctrl *Stats) Events(ctx context.Context, r *request.StatsEvents) (interface{}, error) {
	return ctrl.svc.Events(ctx, r.Kind, service.StatisticsRequest{
		From: r.From,
		To:   r.To,
	})
}
