package service

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza/server/automation/types"
	internalAuth "github.com/cortezaproject/corteza/server/pkg/auth"
	gk "github.com/cortezaproject/corteza/server/pkg/gatekeep"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/modern-go/reflect2"
	"go.uber.org/zap"
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

	gatekeepSvc interface {
		Lock(context.Context, gk.Constraint) (gk.Lock, error)
		Unlock(context.Context, gk.Constraint) error
		ProbeLock(context.Context, gk.Constraint, uint64) (gk.LockState, error)
		ProbeResource(context.Context, string) ([]gk.Lock, error)
		Subscribe(gk.EventListener) int
		Unsubscribe(int)
	}

	gatekeepAccessController interface {
		CanManageGatekeep(context.Context) bool
	}

	gatekeep struct {
		store store.Storer
		ac    gatekeepAccessController
		log   *zap.Logger

		gatekeep   gatekeepSvc
		watcherRef int
		lockSender lockSender
	}

	lockSender interface {
		Send(kind string, payload any, userIDs ...uint64) error
	}
)

func Gatekeep(s store.Storer, log *zap.Logger, ws lockSender) *gatekeep {
	return &gatekeep{
		store:      s,
		ac:         DefaultAccessControl,
		log:        log.Named("gatekeep"),
		gatekeep:   gk.Service(),
		lockSender: ws,
	}
}

func stdLocker() lockerSvc {
	return gk.Locker(gk.Service(), stdLockerFuncs()...)
}

func stdLockerFuncs() (out []gk.LockerConstraint) {
	return []gk.LockerConstraint{
		gk.WithDefaultAwait(),

		// Pull out user ID from context
		//
		// For the sakes of ease, we can use the WorkflowInvoker to avoid issues
		// where we use different users to execute workflows.
		func(ctx context.Context, c gk.Constraint) gk.Constraint {
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

func (svc *gatekeep) Lock(ctx context.Context, c gk.Constraint) (lock gk.Lock, err error) {
	if !svc.ac.CanManageGatekeep(ctx) {
		// @todo proper errors
		err = fmt.Errorf("not allowed to manage gatekeep locks")
		return
	}

	lock, err = svc.gatekeep.Lock(ctx, c)
	if err != nil {
		return
	}

	return
}

func (svc *gatekeep) Unlock(ctx context.Context, c gk.Constraint) (err error) {
	if !svc.ac.CanManageGatekeep(ctx) {
		// @todo proper errors
		err = fmt.Errorf("not allowed to manage gatekeep locks")
		return
	}

	err = svc.gatekeep.Unlock(ctx, c)
	if err != nil {
		return
	}

	return
}

func (svc *gatekeep) Check(ctx context.Context, c gk.Constraint, ref uint64) (state gk.LockState, err error) {
	return svc.gatekeep.ProbeLock(ctx, c, ref)
}

// @todo not handling the unsubscribe at the moment.
// Since this will be unwatched when the server is killed, the listener
// will also go away at that point.
func (svc *gatekeep) Watch(ctx context.Context) {
	gkSvc := svc.gatekeep
	if gkSvc == nil {
		return
	}

	// Subscribe lock/unlock events to dispatch
	svc.watcherRef = gkSvc.Subscribe(func(e gk.Event) {
		switch e.Kind {
		case gk.EbEventLockResolved:
			svc.lockSender.Send("resourceLockResolved", e.Lock)

		case gk.EbEventLockReleased:
			svc.lockSender.Send("resourceLockReleased", e.Lock)
		}
	})
}
