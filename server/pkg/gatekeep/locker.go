package gatekeep

import (
	"context"
	"sync"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/auth"
)

type (
	lockerBit struct {
		Ref        uint64
		Constraint Constraint
	}
	locker struct {
		mux sync.RWMutex

		svc             *service
		locks           []lockerBit
		lockConstraints []lockerConstraint
	}

	lockerConstraint func(c Constraint) Constraint

	identifyable interface {
		Identity() uint64
	}
)

const (
	defaultLockAwait = time.Second * 5
)

func Locker(svc *service, ff ...lockerConstraint) *locker {
	return &locker{
		svc:             svc,
		lockConstraints: ff,
	}
}

// Read attempts to lock the resource for reading
//
// By default, the lock will pend for 5 seconds
func (lg *locker) Read(ctx context.Context, res ...string) (err error) {
	return lg.add(ctx, auth.GetIdentityFromContext(ctx), opRead, res...)
}

// Read attempts to lock the resource for writing
//
// By default, the lock will pend for 5 seconds
func (lg *locker) Write(ctx context.Context, res ...string) (err error) {
	return lg.add(ctx, auth.GetIdentityFromContext(ctx), opWrite, res...)
}

func (lg *locker) add(ctx context.Context, idt identifyable, op Operation, rr ...string) (err error) {
	lg.mux.Lock()
	defer lg.mux.Unlock()

	cc := make([]Constraint, len(rr))
	for i, res := range rr {
		cc[i] = Constraint{
			Resource:  res,
			Operation: op,
			UserID:    idt.Identity(),
		}

		for _, f := range lg.lockConstraints {
			cc[i] = f(cc[i])
		}

		if cc[i].Await == 0 {
			cc[i].Await = defaultLockAwait
		}
	}

	refs, errs := AwaitLocks(ctx, lg.svc, cc...)
	for i, ref := range refs {
		// If one failed, assume all failed
		if errs[i] != nil {
			return errs[i]
		}

		lg.locks = append(lg.locks, lockerBit{
			Ref:        ref,
			Constraint: cc[i],
		})
	}

	return
}

// Free releases all the held locks
func (lg *locker) Free(ctx context.Context) {
	lg.mux.Lock()
	defer lg.mux.Unlock()

	for _, l := range lg.locks {
		// Omit errors
		// @todo should we care about those here?
		lg.svc.Unlock(ctx, l.Constraint)
	}

	return
}
