package gatekeep

import (
	"context"
	"sync"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/auth"
)

type (
	locker struct {
		mux sync.RWMutex

		svc             gksvc
		locks           []lockerBit
		lockConstraints []LockerConstraint
	}

	lockerBit struct {
		Ref        uint64
		Constraint Constraint
	}

	LockerConstraint func(ctx context.Context, c Constraint) Constraint

	gksvc interface {
		Lock(context.Context, Constraint) (Lock, error)
		Unlock(context.Context, Constraint) error
		Subscribe(listener EventListener) int
		Unsubscribe(id int)
	}

	identifyable interface {
		Identity() uint64
	}
)

const (
	defaultLockAwait = time.Second * 5
)

func Locker(svc gksvc, ff ...LockerConstraint) *locker {
	return &locker{
		svc:             svc,
		lockConstraints: ff,
	}
}

func WithDefaultAwait() LockerConstraint {
	return func(_ context.Context, c Constraint) Constraint {
		c.Await = defaultLockAwait
		return c
	}
}

// Read attempts to lock the resource for reading
//
// By default, the lock will pend for 5 seconds
func (lg *locker) Read(ctx context.Context, res ...string) (err error) {
	return lg.add(ctx, OpRead, res...)
}

// Read attempts to lock the resource for writing
//
// By default, the lock will pend for 5 seconds
func (lg *locker) Write(ctx context.Context, res ...string) (err error) {
	return lg.add(ctx, OpWrite, res...)
}

func (lg *locker) add(ctx context.Context, op Operation, rr ...string) (err error) {
	lg.mux.Lock()
	defer lg.mux.Unlock()

	cc := make([]Constraint, len(rr))
	for i, res := range rr {
		cc[i] = Constraint{
			Resource:  res,
			Operation: op,
		}

		for _, f := range lg.lockConstraints {
			cc[i] = f(ctx, cc[i])
		}

		if cc[i].UserID == 0 {
			cc[i].UserID = auth.GetIdentityFromContext(ctx).Identity()
		}
	}

	locks, errs := AwaitLocks(ctx, lg.svc, cc...)
	for i, lock := range locks {
		// If one failed, assume all failed
		if errs[i] != nil {
			return errs[i]
		}

		lg.locks = append(lg.locks, lockerBit{
			Ref:        lock.ID,
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
