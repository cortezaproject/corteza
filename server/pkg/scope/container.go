package scope

import (
	"reflect"
	"sync"
)

// Container holds one instance per interface type for a single scope. The
// any-typed storage is recovered statically via Get[T].
type Container struct {
	mu sync.RWMutex
	m  map[reflect.Type]any
}

func newContainer() *Container {
	return &Container{m: make(map[reflect.Type]any)}
}

func key[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

// Set stores v under the interface type T.
//
// ALWAYS pass the type argument explicitly, naming the interface:
// Set[ExecutionEngine](c, eng). If T is left to inference it resolves to v's
// concrete type, the value is keyed under that concrete type, and a later
// Get[SomeInterface] silently misses.
func Set[T any](c *Container, v T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[key[T]()] = v
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

// MustGet panics if T is absent. Use only where the caller guarantees presence.
func MustGet[T any](c *Container) T {
	v, ok := Get[T](c)
	if !ok {
		panic("scope: component not found: " + key[T]().String())
	}
	return v
}
