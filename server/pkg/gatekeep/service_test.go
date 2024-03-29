package gatekeep

import (
	"context"
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/stretchr/testify/require"
)

func TestBasicLocking(t *testing.T) {
	ctx, svc, req := prep(t)
	var (
		err   error
		ref   uint64
		tt    []lock
		state LockState
	)

	t.Run("acquire read lock", func(t *testing.T) {
		t.Run("User 1 gets read lock", func(t *testing.T) {
			ref, state, err = svc.Lock(ctx, Constraint{
				Resource:  "resource_1/a/b",
				Operation: OpRead,
				UserID:    1,
			})

			req.NoError(err)
			req.Equal(lockStateLocked, state)

			tt, err = svc.probeResource(ctx, "resource_1/a/b")
			req.NoError(err)
			req.Len(tt, 1)
		})

		// User 2 gets same read lock
		t.Run("User 2 gets same read lock", func(t *testing.T) {
			ref, state, err = svc.Lock(ctx, Constraint{
				Resource:  "resource_1/a/b",
				Operation: OpRead,
				UserID:    2,
			})

			req.NoError(err)
			req.Equal(lockStateLocked, state)

			tt, err = svc.probeResource(ctx, "resource_1/a/b")
			req.NoError(err)
			req.Len(tt, 2)
		})

		// User 3 fails write lock
		t.Run("User 3 fails write lock", func(t *testing.T) {
			ref, state, err = svc.Lock(ctx, Constraint{
				Resource:  "resource_1/a/b",
				Operation: OpWrite,
				UserID:    3,
			})

			req.NoError(err)
			req.Equal(lockStateFailed, state)

			tt, err = svc.probeResource(ctx, "resource_1/a/b")
			req.NoError(err)
			req.Len(tt, 2)
		})

		// User 2 release read lock
		t.Run("User 2 release read lock", func(t *testing.T) {
			err = svc.Unlock(ctx, Constraint{
				Resource:  "resource_1/a/b",
				Operation: OpRead,
				UserID:    2,
			})

			req.NoError(err)

			tt, err = svc.probeResource(ctx, "resource_1/a/b")
			req.NoError(err)
			req.Len(tt, 1)
		})

		// User 3 fails write lock
		t.Run("User 3 fails write lock", func(t *testing.T) {
			ref, state, err = svc.Lock(ctx, Constraint{
				Resource:  "resource_1/a/b",
				Operation: OpWrite,
				UserID:    3,
			})

			req.NoError(err)
			req.Equal(lockStateFailed, state)

			tt, err = svc.probeResource(ctx, "resource_1/a/b")
			req.NoError(err)
			req.Len(tt, 1)
		})

		// User 1 release read lock
		t.Run("User 1 release read lock", func(t *testing.T) {
			err = svc.Unlock(ctx, Constraint{
				Resource:  "resource_1/a/b",
				Operation: OpRead,
				UserID:    1,
			})

			req.NoError(err)

			tt, err = svc.probeResource(ctx, "resource_1/a/b")
			req.NoError(err)
			req.Len(tt, 0)
		})

		// User 3 gets write lock
		t.Run("User 3 gets write lock", func(t *testing.T) {
			ref, state, err = svc.Lock(ctx, Constraint{
				Resource:  "resource_1/a/b",
				Operation: OpWrite,
				UserID:    3,
			})

			req.NoError(err)
			req.Equal(lockStateLocked, state)

			tt, err = svc.probeResource(ctx, "resource_1/a/b")
			req.NoError(err)
			req.Len(tt, 1)
		})

		t.Run("User 4 fails read lock", func(t *testing.T) {
			ref, state, err = svc.Lock(ctx, Constraint{
				Resource:  "resource_1/a/b",
				Operation: OpRead,
				UserID:    4,
			})

			req.NoError(err)
			req.Equal(lockStateFailed, state)

			tt, err = svc.probeResource(ctx, "resource_1/a/b")
			req.NoError(err)
			req.Len(tt, 1)
		})

		t.Run("User 5 fails write lock", func(t *testing.T) {
			ref, state, err = svc.Lock(ctx, Constraint{
				Resource:  "resource_1/a/b",
				Operation: OpWrite,
				UserID:    5,
			})

			req.NoError(err)
			req.Equal(lockStateFailed, state)

			tt, err = svc.probeResource(ctx, "resource_1/a/b")
			req.NoError(err)
			req.Len(tt, 1)
		})
	})

	_ = ref
	_ = state
}

func TestQueueing(t *testing.T) {
	ctx, svc, req := prep(t)
	var (
		err   error
		ref   uint64
		tt    []lock
		state LockState
	)

	_, _, err = svc.Lock(ctx, Constraint{
		Resource:  "resource_1/a/b",
		Operation: OpWrite,
		UserID:    1,
	})
	req.NoError(err)

	var refLockR1 uint64
	t.Run("queue read lock 1", func(t *testing.T) {
		c := Constraint{
			Resource:  "resource_1/a/b",
			Operation: OpRead,
			UserID:    2,
			Await:     time.Hour * 2,
		}
		refLockR1, state, err = svc.Lock(ctx, c)

		req.NoError(err)
		req.Equal(lockStateQueued, state)

		tt, err = svc.probeResource(ctx, "resource_1/a/b")
		req.NoError(err)
		req.Len(tt, 2)

		state, err = svc.ProbeLock(ctx, c, refLockR1)
		req.NoError(err)
		req.Equal(lockStateQueued, state)
	})

	var refLockR2 uint64
	t.Run("queue read lock 2", func(t *testing.T) {
		c := Constraint{
			Resource:  "resource_1/a/b",
			Operation: OpRead,
			UserID:    3,
			Await:     time.Hour * 2,
		}
		refLockR2, state, err = svc.Lock(ctx, c)

		req.NoError(err)
		req.Equal(lockStateQueued, state)

		tt, err = svc.probeResource(ctx, "resource_1/a/b")
		req.NoError(err)
		req.Len(tt, 3)

		state, err = svc.ProbeLock(ctx, c, refLockR2)
		req.NoError(err)
		req.Equal(lockStateQueued, state)
	})

	var refLockW2 uint64
	t.Run("queue write lock 2", func(t *testing.T) {
		c := Constraint{
			Resource:  "resource_1/a/b",
			Operation: OpWrite,
			UserID:    4,
			Await:     time.Hour * 2,
		}
		refLockW2, state, err = svc.Lock(ctx, c)

		req.NoError(err)
		req.Equal(lockStateQueued, state)

		tt, err = svc.probeResource(ctx, "resource_1/a/b")
		req.NoError(err)
		req.Len(tt, 4)

		state, err = svc.ProbeLock(ctx, c, refLockW2)
		req.NoError(err)
		req.Equal(lockStateQueued, state)
	})

	t.Run("release write lock 1", func(t *testing.T) {
		c := Constraint{
			Resource:  "resource_1/a/b",
			Operation: OpWrite,
			UserID:    1,
		}
		err = svc.Unlock(ctx, c)
		req.NoError(err)

		tt, err = svc.probeResource(ctx, "resource_1/a/b")
		req.NoError(err)
		req.Len(tt, 3)

		state, err = svc.ProbeLock(ctx, c, refLockR1)
		req.NoError(err)
		req.Equal(lockStateLocked, state)

		state, err = svc.ProbeLock(ctx, c, refLockR2)
		req.NoError(err)
		req.Equal(lockStateLocked, state)
	})

	t.Run("release write lock 1", func(t *testing.T) {
		c := Constraint{
			Resource:  "resource_1/a/b",
			Operation: OpWrite,
			UserID:    1,
		}
		err = svc.Unlock(ctx, c)
		req.NoError(err)

		tt, err = svc.probeResource(ctx, "resource_1/a/b")
		req.NoError(err)
		req.Len(tt, 3)

		state, err = svc.ProbeLock(ctx, c, refLockR1)
		req.NoError(err)
		req.Equal(lockStateLocked, state)

		state, err = svc.ProbeLock(ctx, c, refLockR2)
		req.NoError(err)
		req.Equal(lockStateLocked, state)
	})

	t.Run("release read lock 1", func(t *testing.T) {
		c := Constraint{
			Resource:  "resource_1/a/b",
			Operation: OpRead,
			UserID:    2,
		}
		err = svc.Unlock(ctx, c)
		req.NoError(err)

		tt, err = svc.probeResource(ctx, "resource_1/a/b")
		req.NoError(err)
		req.Len(tt, 2)

		state, err = svc.ProbeLock(ctx, c, refLockR2)
		req.NoError(err)
		req.Equal(lockStateLocked, state)

		state, err = svc.ProbeLock(ctx, c, refLockW2)
		req.NoError(err)
		req.Equal(lockStateQueued, state)
	})

	t.Run("release read lock 2", func(t *testing.T) {
		c := Constraint{
			Resource:  "resource_1/a/b",
			Operation: OpRead,
			UserID:    3,
		}
		err = svc.Unlock(ctx, c)
		req.NoError(err)

		tt, err = svc.probeResource(ctx, "resource_1/a/b")
		req.NoError(err)
		req.Len(tt, 1)

		state, err = svc.ProbeLock(ctx, c, refLockW2)
		req.NoError(err)
		req.Equal(lockStateLocked, state)
	})

	_ = ref
	_ = tt
}

func TestResourceHierarchy(t *testing.T) {
	ctx, svc, req := prep(t)
	var (
		err   error
		ref   uint64
		state LockState
	)

	t.Run("acquire same kind lock", func(t *testing.T) {
		t.Run("user 1 acquire generic write lock", func(t *testing.T) {
			ref, state, err = svc.Lock(ctx, Constraint{
				Resource:  "resource_1/a/*",
				Operation: OpWrite,
				UserID:    1,
			})

			req.NoError(err)
			req.Equal(lockStateLocked, state)
		})

		t.Run("user 2 fails specific write lock", func(t *testing.T) {
			ref, state, err = svc.Lock(ctx, Constraint{
				Resource:  "resource_1/a/b",
				Operation: OpWrite,
				UserID:    2,
			})

			req.NoError(err)
			req.Equal(lockStateFailed, state)
		})

		// Cleanup
		svc.Unlock(ctx, Constraint{
			Resource:  "resource_1/a/*",
			Operation: OpWrite,
			UserID:    1,
		})
	})

	t.Run("acquire different kind lock", func(t *testing.T) {
		t.Run("user 1 acquire generic read lock", func(t *testing.T) {
			ref, state, err = svc.Lock(ctx, Constraint{
				Resource:  "resource_1/a/*",
				Operation: OpRead,
				UserID:    1,
			})

			req.NoError(err)
			req.Equal(lockStateLocked, state)
		})

		t.Run("user 2 fails specific write lock", func(t *testing.T) {
			ref, state, err = svc.Lock(ctx, Constraint{
				Resource:  "resource_1/a/b",
				Operation: OpWrite,
				UserID:    2,
			})

			req.NoError(err)
			req.Equal(lockStateFailed, state)
		})

		// Cleanup
		svc.Lock(ctx, Constraint{
			Resource:  "resource_1/a/*",
			Operation: OpRead,
			UserID:    1,
		})
	})

	_ = ref
	_ = state
}

func prep(t require.TestingT) (ctx context.Context, svc *service, req *require.Assertions) {
	req = require.New(t)
	ctx = context.Background()

	id.Init(ctx)

	svc, err := New(nil, InmemStore())
	req.NoError(err)

	return
}
