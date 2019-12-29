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

		// Type of event fired
		EventType() string

		// Match tests if given constraints match
		// event's internal values
		Match(name string, op string, values ...string) bool
	}

	eventbus struct {
		l *sync.RWMutex

		// list of registered handlers
		triggers map[uintptr]*trigger
	}
)

var (
	// Global eventbus
	gEventBus *eventbus
)

func init() {
	gEventBus = New()
}

// Returns
func Default() *eventbus {
	return gEventBus
}

func New() *eventbus {
	return &eventbus{
		l:        &sync.RWMutex{},
		triggers: make(map[uintptr]*trigger),
	}
}

// WaitFor is synchronous event dispatcher
//
// It waits for all handlers and fails on first error
func (b *eventbus) WaitFor(ctx context.Context, ev Event) (err error) {
	b.l.RLock()
	defer b.l.RUnlock()
	for _, t := range b.find(ev) {
		if err = t.Handle(ctx, ev); err != nil {
			return err
		}
	}

	return
}

// Dispatch runs events asynchronously
func (b *eventbus) Dispatch(ctx context.Context, ev Event) {
	b.l.RLock()
	defer b.l.RUnlock()
	for _, t := range b.find(ev) {
		go func(ctx context.Context, t *trigger) {
			_ = t.Handle(ctx, ev)
		}(ctx, t)
	}
}

// Finds all registered triggers compatible with given event
//
// It returns sorted
//
// There is still room for improvement (performance wise) by indexing
// resources and events of each trigger.
func (b *eventbus) find(ev Event) (tt TriggerSet) {
	for _, t := range b.triggers {
		if ev == nil {
			continue
		}

		if !t.Match(ev) {
			continue
		}

		tt = append(tt, t)
	}

	sort.Sort(tt)

	return
}

// Register creates a new trigger with given handler, resource, event with other options and constraints
//
// It returns a trigger identifier that can be used to remove (unregister) trigger later
func (b *eventbus) Register(h Handler, ops ...TriggerRegOp) uintptr {
	b.l.Lock()
	defer b.l.Unlock()

	var (
		trigger = NewTrigger(h, ops...)
		ptr     = uintptr(unsafe.Pointer(trigger))
	)

	b.triggers[ptr] = trigger
	return ptr
}

// Unregister removes one or more registered triggers
func (b *eventbus) Unregister(ptrs ...uintptr) {
	b.l.Lock()
	defer b.l.Unlock()

	for _, ptr := range ptrs {
		delete(b.triggers, ptr)
	}
}
