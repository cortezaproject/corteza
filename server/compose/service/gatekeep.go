package service

import (
	"context"

	"github.com/cortezaproject/corteza/server/automation/types"
	internalAuth "github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/gatekeep"
	"github.com/modern-go/reflect2"
)

type (
	identifyable interface {
		Identity() uint64
	}

	lockerSvc interface {
		Read(context.Context, ...string) error
		Write(context.Context, ...string) error
		Free(context.Context)
	}
)

func stdLocker() lockerSvc {
	return gatekeep.Locker(gatekeep.Service(), stdLockerFuncs()...)
}

func stdLockerFuncs() (out []gatekeep.LockerConstraint) {
	return []gatekeep.LockerConstraint{
		gatekeep.WithDefaultAwait(),

		// Pull out user ID from context
		//
		// For the sakes of ease, we can use the WorkflowInvoker to avoid issues
		// where we use different users to execute workflows.
		func(ctx context.Context, c gatekeep.Constraint) gatekeep.Constraint {
			var id uint64

			aux := ctx.Value(types.WorkflowInvokerCtxKey{})
			if !reflect2.IsNil(aux) {
				id = aux.(identifyable).Identity()
			} else {
				id = internalAuth.GetIdentityFromContext(ctx).Identity()
			}

			c.UserID = id

			return c
		},
	}
}
