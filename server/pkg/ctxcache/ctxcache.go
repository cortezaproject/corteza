package ctxcache

// Keeps values loaded during one unit of work (a workflow run) in its
// context, so a lookup repeated inside that work reads the database once
//
// Entries expire after a short time; a long-running or resumed run picks up
// changes made in the meantime.

import (
	"context"
	"sync"
	"time"
)

type (
	Cache struct {
		ttl     time.Duration
		mux     sync.Mutex
		entries map[string]entry
	}

	entry struct {
		val     any
		expires time.Time
	}

	ctxKey struct{}
)

var now = time.Now

func New(ttl time.Duration) *Cache {
	return &Cache{ttl: ttl, entries: make(map[string]entry)}
}

// With returns a context carrying the given cache
func With(ctx context.Context, c *Cache) context.Context {
	return context.WithValue(ctx, ctxKey{}, c)
}

// From returns the cache the context carries, nil when there is none
func From(ctx context.Context) *Cache {
	c, _ := ctx.Value(ctxKey{}).(*Cache)
	return c
}

// Get returns the value cached under the key or loads and caches it
//
// Without a cache in the context it always loads. Errors are not cached.
func Get[T any](ctx context.Context, key string, load func() (T, error)) (T, error) {
	c := From(ctx)
	if c == nil {
		return load()
	}

	c.mux.Lock()
	e, ok := c.entries[key]
	c.mux.Unlock()

	if ok && now().Before(e.expires) {
		return e.val.(T), nil
	}

	val, err := load()
	if err != nil {
		return val, err
	}

	c.mux.Lock()
	c.entries[key] = entry{val: val, expires: now().Add(c.ttl)}
	c.mux.Unlock()

	return val, nil
}
