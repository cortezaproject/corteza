package service

import (
	"context"

	"github.com/cortezaproject/corteza/server/store"

	"github.com/cortezaproject/corteza/server/pkg/actionlog"
	"github.com/cortezaproject/corteza/server/pkg/rbac"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	statistics struct {
		ac        statsAccessControl
		actionlog actionlog.Recorder
		store     store.Storer
	}

	statsAccessControl interface {
		CanSearchUsers(context.Context) bool
		CanSearchRoles(context.Context) bool
		CanSearchApplications(context.Context) bool
	}

	StatisticsMetricsPayload struct {
		Users        *types.UserMetrics        `json:"users"`
		Roles        *types.RoleMetrics        `json:"roles"`
		Applications *types.ApplicationMetrics `json:"applications"`

		Rbac rbac.Stats `json:"rbac"`
	}
)

func Statistics() *statistics {
	return &statistics{
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
	}
}

// Metrics collects relevant metrics and returns it
//
// @todo remove this service and move it to rest ctrl layer
func (svc statistics) Metrics(ctx context.Context) (rval *StatisticsMetricsPayload, err error) {
	err = func() error {
		rval = &StatisticsMetricsPayload{}

		if svc.ac.CanSearchUsers(ctx) {
			if rval.Users, err = store.UserMetrics(ctx, svc.store); err != nil {
				return err
			}
		}

		if svc.ac.CanSearchRoles(ctx) {
			if rval.Roles, err = store.RoleMetrics(ctx, svc.store); err != nil {
				return err
			}
		}

		if svc.ac.CanSearchApplications(ctx) {
			if rval.Applications, err = store.ApplicationMetrics(ctx, svc.store); err != nil {
				return err
			}
		}

		// @todo rbac
		rval.Rbac, err = rbac.Global().Stats()
		if err != nil {
			return err
		}

		return nil
	}()

	return rval, svc.recordAction(ctx, &statisticsActionProps{}, StatisticsActionServe, err)
}
