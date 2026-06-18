package scope

import (
	"context"
	"errors"
	"reflect"
	"sync"
)

// Teardownable is implemented by components that hold resources (goroutines,
// pools, watchers) needing release when their scope is evicted.
type Teardownable interface {
	Close(ctx context.Context) error
}

// Container holds one instance per interface type for a single scope.
// The any-typed storage is recovered statically via Get[T].
type Container struct {
	mu    sync.RWMutex
	m     map[reflect.Type]any
	order []reflect.Type // insertion order; teardown walks it in reverse
}

func newContainer() *Container {
	return &Container{m: make(map[reflect.Type]any)}
}

func key[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

// Set stores v under the interface type T. Re-setting T replaces in place
// (does not duplicate the order entry).
//
// ALWAYS pass the type argument explicitly, naming the interface:
// Set[ExecutionEngine](c, eng). If T is left to inference it resolves to v's
// concrete type, the value is keyed under that concrete type, and a later
// Get[SomeInterface] silently misses.
func Set[T any](c *Container, v T) {
	k := key[T]()
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, exists := c.m[k]; !exists {
		c.order = append(c.order, k)
	}
	c.m[k] = v
}

// Get returns the component stored under interface T.
func Get[T any](c *Container) (T, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.m[key[T]()]
	if !ok {
		var zero T
		return zero, false
	}
	return v.(T), true
}

// MustGet panics if T is absent. Use only where a provider guarantees presence.
func MustGet[T any](c *Container) T {
	v, ok := Get[T](c)
	if !ok {
		panic("scope: component not found: " + key[T]().String())
	}
	return v
}

// close tears down every Teardownable in reverse insertion order, joining
// errors. It snapshots and clears the container under the lock, then runs
// teardowns OUTSIDE it: a component's Close may safely call Get on the same
// container (it observes an already-empty container instead of deadlocking
// on the non-reentrant RWMutex).
func (c *Container) close(ctx context.Context) error {
	c.mu.Lock()
	snapshot := make([]any, 0, len(c.order))
	for i := len(c.order) - 1; i >= 0; i-- { // reverse insertion order
		snapshot = append(snapshot, c.m[c.order[i]])
	}
	c.m = make(map[reflect.Type]any)
	c.order = nil
	c.mu.Unlock()

	var errs []error
	for _, v := range snapshot {
		if td, ok := v.(Teardownable); ok {
			if err := td.Close(ctx); err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errors.Join(errs...)
}
