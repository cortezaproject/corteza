package service

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/system/repository"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	statistics struct {
		actionlog actionlog.Recorder
		ac        statisticsAccessController
	}

	statisticsAccessController interface {
		CanAccess(context.Context) bool
	}

	StatisticsMetricsPayload struct {
		Users        *types.UserMetrics        `json:"users"`
		Roles        *types.RoleMetrics        `json:"roles"`
		Applications *types.ApplicationMetrics `json:"applications"`
	}
)

func Statistics() *statistics {
	return &statistics{
		actionlog: DefaultActionlog,
		ac:        DefaultAccessControl,
	}
}

func (svc statistics) Metrics(ctx context.Context) (rval *StatisticsMetricsPayload, err error) {
	db := repository.DB(ctx)

	err = db.Transaction(func() error {
		if !svc.ac.CanAccess(ctx) {
			return StatisticsErrNotAllowedToReadStatistics()
		}

		rval = &StatisticsMetricsPayload{}

		if rval.Users, err = repository.User(ctx, db).Metrics(); err != nil {
			return err
		}

		if rval.Roles, err = repository.Role(ctx, db).Metrics(); err != nil {
			return err
		}

		if rval.Applications, err = repository.Application(ctx, db).Metrics(); err != nil {
			return err
		}

		return nil
	})

	return rval, svc.recordAction(ctx, &statisticsActionProps{}, StatisticsActionServe, err)
}
