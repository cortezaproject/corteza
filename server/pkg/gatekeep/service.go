package gatekeep

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/davecgh/go-spew/spew"
	"go.uber.org/zap"
)

type (
	service struct {
		mux sync.RWMutex

		store        store
		queueManager *queueManager

		events eventManager

		log *zap.Logger
	}

	eventListener func(evt Event)
	Event         struct {
		Kind ebEvent
		Lock lock
	}

	eventManager interface {
		Subscribe(listener eventListener) int
		Unsubscribe(int)
		Publish(event Event)
	}

	Constraint struct {
		id uint64

		Resource  string
		Operation Operation
		UserID    uint64
		Overwrite bool
		Await     time.Duration

		queuedAt time.Time
	}

	queue struct {
		queue []Constraint
	}

	queueManager struct {
		mux    sync.Mutex
		queues map[string]*queue
	}

	store interface {
		GetValue(ctx context.Context, key string) ([]byte, error)
		SetValue(ctx context.Context, key string, v []byte) error
		DeleteValue(ctx context.Context, key string) error
	}

	lock struct {
		ID        uint64    `json:"id"`
		UserID    uint64    `json:"userID"`
		CreatedAt time.Time `json:"createdAt"`
		Resource  string    `json:"resource"`
		Operation Operation `json:"operation"`

		State lockState `json:"state"`

		LockDuration time.Duration `json:"lockDuration"`
		LockExpires  *time.Time    `json:"lockExpires"`
	}

	ebEvent int

	lockState int
	Operation string
)

const (
	opRead  Operation = "read"
	opWrite Operation = "write"
)

const (
	lockStateNil lockState = iota
	lockStateLocked
	lockStateFailed
	lockStateQueued
)

const (
	ebEventLockResolved ebEvent = iota
	ebEventLockReleased
)

var (
	gSvc *service

	// wrapper around id.Next() that will aid service testing
	nextID = func() uint64 {
		return id.Next()
	}
)

// New creates a DAL service with the primary connection
//
// It needs an established and working connection to the primary store
func New(log *zap.Logger, s store) (*service, error) {
	svc := &service{
		mux:   sync.RWMutex{},
		log:   log,
		store: s,

		queueManager: &queueManager{
			mux:    sync.Mutex{},
			queues: make(map[string]*queue),
		},

		events: &inMemBuss{},
	}
	return svc, nil
}

func Initialized() bool {
	return gSvc != nil
}

// Service returns the global initialized DAL service
//
// Function will panic if DAL service is not set (via SetGlobal)
func Service() *service {
	if gSvc == nil {
		panic("gatekeep global service not initialized: call gatekeep.SetGlobal first")
	}

	return gSvc
}

func SetGlobal(svc *service, err error) {
	if err != nil {
		panic(err)
	}

	gSvc = svc
}

// Lock attempts to acquire a lock conforming to the given constraints
//
// If a lock can't be acquired the request will either be queued or fail
// (if the .Await field is not set)
//
// The function doesn't block/wait for the lock to be acquired; that needs
// to be done by the caller
func (svc *service) Lock(ctx context.Context, c Constraint) (ref uint64, state lockState, err error) {
	svc.mux.Lock()
	defer svc.mux.Unlock()

	// Probe existing resource locks so we can figure out what we can do
	ll, err := svc.probeResource(ctx, c.Resource)
	if err != nil {
		return
	}

	// Check if we already have this lock so we can potentially extend the lock
	for _, l := range ll {
		if l.matchesConstraints(c) {
			// @todo extending?
			// @todo queued
			return l.ID, l.State, nil
		}
	}

	// If we're wanting to acquire a read lock, we can only of there are none
	// or all existing locks are also read locks
	allRead := c.Operation == opRead
	for _, t := range ll {
		allRead = allRead && t.Operation == opRead
	}

	// If there are locks and we're not willing to wait, we're done
	if (len(ll) > 0 && !allRead) && c.Await == 0 {
		state = lockStateFailed
		return
	}

	// If there are no locks or all are read locks, we can acquire the lock
	if len(ll) == 0 || allRead {
		ref, err = svc.acquireLock(ctx, c)
		if err != nil {
			state = lockStateFailed
			return
		}

		state = lockStateLocked
		return
	}

	// Queue the lock
	ref, err = svc.queueManager.queueLock(ctx, c)
	state = lockStateQueued
	return
}

// Unlock releases the lock or unqueues the lock if it's queued
//
// The function won't error out if the lock doesn't exist
// @todo should it?
func (svc *service) Unlock(ctx context.Context, c Constraint) (err error) {
	svc.mux.Lock()
	defer svc.mux.Unlock()
	// releasing a lock may result in other locks being acquirable
	defer svc.doQueued(ctx, c)

	ref, exists, err := svc.check(ctx, c)
	if err != nil {
		return
	}

	if ref == 0 {
		return
	}

	if exists == lockStateLocked {
		return svc.releaseLock(ctx, c, ref)
	} else if exists == lockStateQueued {
		return svc.releaseQueued(ctx, c, ref)
	}

	return
}

// ProbeLock returns the current state of the lock
func (svc *service) ProbeLock(ctx context.Context, c Constraint, ref uint64) (state lockState, err error) {
	svc.mux.Lock()
	defer svc.mux.Unlock()

	tt, err := svc.probeResource(ctx, c.Resource)
	if err != nil {
		return
	}

	for _, t := range tt {
		if t.ID == ref {
			return t.State, nil
		}
	}

	return
}

func (svc *service) ProbeResource(ctx context.Context, r string) (tt []lock, err error) {
	svc.mux.RLock()
	defer svc.mux.RUnlock()

	return svc.probeResource(ctx, r)
}

// probeResource returns all of the locks on the given resource
//
// The function returns both already acquired and queued locks
func (svc *service) probeResource(ctx context.Context, r string) (tt []lock, err error) {
	// Get the currently stored locks
	bb, err := svc.store.GetValue(ctx, r)
	if err != nil && err.Error() == "not found" {
		return tt, nil
	}
	if err != nil {
		return
	}

	err = json.Unmarshal(bb, &tt)
	if err != nil {
		return
	}

	// Get queued locks
	aux := svc.queueManager.queues[r]
	if aux == nil {
		return
	}

	for _, c := range aux.queue {
		tt = append(tt, lock{
			ID:        c.id,
			UserID:    c.UserID,
			Resource:  c.Resource,
			Operation: c.Operation,
			State:     lockStateQueued,
		})
	}

	// @todo
	return
}

// check returns the lock reference along with it's state
func (svc *service) check(ctx context.Context, c Constraint) (ref uint64, state lockState, err error) {
	aux, err := svc.probeResource(ctx, c.Resource)
	if err != nil {
		return
	}

	for _, t := range aux {
		if !t.matchesConstraints(c) {
			continue
		}

		return t.ID, t.State, nil
	}

	return 0, lockStateNil, nil
}

// Cleanup will release stale things
func (svc *service) Cleanup(ctx context.Context) (err error) {
	svc.mux.Lock()
	defer svc.mux.Unlock()

	// @todo cleanup stale/overdue locks

	// Cleanup stale queued locks
	qm := svc.queueManager
	if qm == nil {
		return
	}

	qm.mux.Lock()
	defer qm.mux.Unlock()

	now := time.Now()
	// Go backwards and spice out the ones that need to be removed.
	// Broadcast down the buss so we can kill off the watchers.
	for _, qq := range qm.queues {
		for i := len(qq.queue) - 1; i >= 0; i-- {
			c := qq.queue[i]
			l := lock{
				ID:        c.id,
				UserID:    c.UserID,
				CreatedAt: c.queuedAt,
				Resource:  c.Resource,
				Operation: c.Operation,

				State: lockStateFailed,
			}

			if c.queuedAt.IsZero() {
				qq.queue = append(qq.queue[:i], qq.queue[i+1:]...)

				svc.events.Publish(Event{
					Kind: ebEventLockResolved,
					Lock: l,
				})
				continue
			}

			if now.After(c.queuedAt.Add(c.Await)) {
				qq.queue = append(qq.queue[:i], qq.queue[i+1:]...)

				svc.events.Publish(Event{
					Kind: ebEventLockResolved,
					Lock: l,
				})
				continue
			}
		}
	}

	return
}

func (svc *service) Watch(ctx context.Context) {
	tCleanup := time.NewTicker(5 * time.Second)
	tQueued := time.NewTicker(1 * time.Second)

	var err error
	go func() {
		for {
			select {
			case <-tQueued.C:
				// @todo

			case <-tCleanup.C:
				err = svc.Cleanup(ctx)
				if err != nil {
					// @todo logging
					spew.Dump("cleanup error", err)
					err = nil
				}

			case <-ctx.Done():
				tCleanup.Stop()
				return
			}
		}
	}()
}

// @todo we could consider prioritizing some/all read locks over write locks
// so we can have a higher throughput
func (qm *queueManager) queueLock(ctx context.Context, c Constraint) (ref uint64, err error) {
	qm.mux.Lock()
	defer qm.mux.Unlock()

	key := c.Resource

	_, ok := qm.queues[key]
	if !ok {
		qm.queues[key] = &queue{
			queue: make([]Constraint, 0, 24),
		}
	}

	q := qm.queues[key]
	c.id = nextID()
	c.queuedAt = time.Now()
	q.queue = append(q.queue, c)
	qm.queues[key] = q

	return c.id, nil
}

func (svc *service) doQueued(ctx context.Context, c Constraint) (err error) {
	svc.queueManager.mux.Lock()
	defer svc.queueManager.mux.Unlock()

	q := svc.queueManager.queues[c.Resource]
	if q == nil {
		return
	}

	if len(q.queue) == 0 {
		delete(svc.queueManager.queues, c.Resource)
		return
	}

	doReads := q.queue[0].Operation == opRead

	if !doReads {
		// Check if we can acquire a new one
		qc := q.queue[0]

		// Probe existing resource locks so we can figure out what we can do
		var tt []lock
		tt, err = svc.probeResource(ctx, qc.Resource)
		if err != nil {
			return
		}

		// Check if we already have this lock so we can potentially extend the lock
		for _, t := range tt {
			if t.ID == qc.id {
				continue
			}

			// If there are any locks and we're trying a write lock; no bueno
			return
		}

		q.queue = q.queue[1:]

		// @todo
		_, err = svc.acquireLock(ctx, qc, qc.id)
		if err != nil {
			// @todo logging
			spew.Dump("doQueued error", err)
		}

		return
	}

	var i int
	var qc Constraint
	for i, qc = range q.queue {
		if qc.Operation != opRead {
			break
		}

		_, err = svc.acquireLock(ctx, qc, qc.id)
		if err != nil {
			// @todo logging
			spew.Dump("doQueued error", err)
			continue
		}
	}

	q.queue = q.queue[i:]
	return
}

func (svc *service) acquireLock(ctx context.Context, c Constraint, ids ...uint64) (ref uint64, err error) {
	tt := make([]lock, 0)

	// Get current locks from the store
	// @todo we can probably pass the OG slice around
	baseB, err := svc.store.GetValue(ctx, c.Resource)
	if err != nil && err.Error() != "not found" {
		return
	}

	if len(baseB) > 0 {
		err = json.Unmarshal(baseB, &tt)
		if err != nil {
			return
		}
	}

	id := nextID()
	if len(ids) > 0 {
		id = ids[0]
	}
	tt = append(tt, lock{
		ID:        id,
		UserID:    c.UserID,
		CreatedAt: time.Now(),
		Resource:  c.Resource,
		Operation: c.Operation,
		State:     lockStateLocked,
	})

	bb, err := json.Marshal(tt)
	if err != nil {
		return
	}

	err = svc.store.SetValue(ctx, c.Resource, bb)
	if err != nil {
		return
	}

	ref = id
	return
}

// releaseLock removes the lock from the store
func (svc *service) releaseLock(ctx context.Context, c Constraint, ref uint64) (err error) {
	baseB, err := svc.store.GetValue(ctx, c.Resource)
	if err != nil && err.Error() != "not found" {
		return
	}

	tt := make([]lock, 0)
	if len(baseB) > 0 {
		err = json.Unmarshal(baseB, &tt)
		if err != nil {
			return
		}
	}

	aux := make([]lock, 0, len(tt))
	for _, t := range tt {
		if t.ID == ref {
			continue
		}
		aux = append(aux, t)
	}

	bb, err := json.Marshal(aux)
	if err != nil {
		return
	}

	return svc.store.SetValue(ctx, c.Resource, bb)
}

// releaseQueued removes the lock from the queue
func (svc *service) releaseQueued(ctx context.Context, c Constraint, ref uint64) (err error) {
	if svc.queueManager.queues == nil {
		return
	}

	svc.queueManager.mux.Lock()
	defer svc.queueManager.mux.Unlock()

	qq := svc.queueManager.queues[c.Resource]
	if qq == nil {
		return
	}

	aux := make([]Constraint, 0, 24)
	for _, q := range qq.queue {
		if q.id == ref {
			continue
		}
		aux = append(aux, q)
	}

	if len(aux) == 0 {
		delete(svc.queueManager.queues, c.Resource)
		return
	}

	svc.queueManager.queues[c.Resource].queue = aux
	return
}

func (t lock) matchesConstraints(c Constraint) (ok bool) {
	return t.UserID == c.UserID && t.Resource == c.Resource && t.Operation == c.Operation
}
