package eventbus

import (
	"context"
	"sort"
	"sync"
	"unsafe"
)

type (
	Event interface {
		// ResourceType from resource that fired the event
		ResourceType() string

		// EventType returns type of event fired
		EventType() string

		// Match tests if given constraints match
		// event's internal values
		Match(ConstraintMatcher) bool
	}

	eventbus struct {
		// waitgroup for dispatch
		wg *sync.WaitGroup

		// Read & write locking
		// prevent event handling during handler (un)registration
		l *sync.RWMutex

		// list of registered handlers
		handlers map[uintptr]*handler
	}
)

var (
	// Global eventbus
	gEventBus *eventbus
)

func init() {
	gEventBus = New()
}

// Service returns global event bus service
func Service() *eventbus {
	return gEventBus
}

func Set(eb *eventbus) {
	gEventBus = eb
}

func New() *eventbus {
	return &eventbus{
		wg:       &sync.WaitGroup{},
		l:        &sync.RWMutex{},
		handlers: make(map[uintptr]*handler),
	}
}

// WaitFor is synchronous event dispatcher
//
// It waits for each handler and fails on first error
func (b *eventbus) WaitFor(ctx context.Context, ev Event) (err error) {
	b.l.RLock()
	defer b.l.RUnlock()

	for _, t := range b.find(ev) {
		err = func(ctx context.Context, t *handler) error {
			b.wg.Add(1)
			defer b.wg.Done()
			return t.Handle(ctx, ev)

		}(ctx, t)

		if err != nil {
			return
		}
	}

	return
}

// Dispatch runs events asynchronously
func (b *eventbus) Dispatch(ctx context.Context, ev Event) {
	b.l.RLock()
	defer b.l.RUnlock()
	for _, t := range b.find(ev) {
		b.wg.Add(1)
		go func(ctx context.Context, t *handler) {
			defer b.wg.Done()
			_ = t.Handle(ctx, ev)
		}(ctx, t)
	}
}

// Waits for all dispatched events
//
// Should only be used for testing
func (b *eventbus) wait() {
	b.wg.Wait()
}

// Finds all registered handlers compatible with given event
//
// It returns sorted handlers
//
// There is still room for improvement (performance wise) by indexing
// resources and events of each handler.
func (b *eventbus) find(ev Event) (tt HandlerSet) {
	if ev == nil {
		return
	}

	for _, t := range b.handlers {

		if !t.Match(ev) {
			continue
		}

		tt = append(tt, t)
	}

	sort.Sort(tt)

	return
}

// Register creates a new handler with given handler, resource, event with other options and constraints
//
// It returns a handler identifier that can be used to remove (unregister) handler later
func (b *eventbus) Register(h HandlerFn, ops ...HandlerRegOp) uintptr {
	b.l.Lock()
	defer b.l.Unlock()

	var (
		handlers = NewHandler(h, ops...)
		ptr      = uintptr(unsafe.Pointer(handlers))
	)

	b.handlers[ptr] = handlers
	return ptr
}

// Unregister removes one or more registered handlers
func (b *eventbus) Unregister(ptrs ...uintptr) {
	b.l.Lock()
	defer b.l.Unlock()

	for _, ptr := range ptrs {
		delete(b.handlers, ptr)
	}
}

// UnregisterByResource removes one or more registered handlers that match the given resource
func (b *eventbus) UnregisterByResource(r string) {
	b.l.Lock()
	defer b.l.Unlock()
	for p, h := range b.handlers {
		if h.resourceTypes[r] {
			delete(b.handlers, p)
			continue
		}

	}
}
