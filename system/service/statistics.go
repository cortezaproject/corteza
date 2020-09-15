package service

import (
	"context"
	"github.com/cortezaproject/corteza-server/store"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	statistics struct {
		actionlog actionlog.Recorder
		ac        statisticsAccessController
		store     store.Storer
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
		store:     DefaultStore,
	}
}

// Metrics collects relevant metrics and returns it
//
// @todo remove this service and move it to rest ctrl layer
func (svc statistics) Metrics(ctx context.Context) (rval *StatisticsMetricsPayload, err error) {
	err = func() error {
		if !svc.ac.CanAccess(ctx) {
			return StatisticsErrNotAllowedToReadStatistics()
		}

		rval = &StatisticsMetricsPayload{}

		if rval.Users, err = svc.store.UserMetrics(ctx); err != nil {
			return err
		}

		if rval.Roles, err = svc.store.RoleMetrics(ctx); err != nil {
			return err
		}

		if rval.Applications, err = svc.store.ApplicationMetrics(ctx); err != nil {
			return err
		}

		return nil
	}()

	return rval, svc.recordAction(ctx, &statisticsActionProps{}, StatisticsActionServe, err)
}
