package rbac

import "sync"

const userCacheMaxSize = 1024

// userCache is a bounded, goroutine-safe cache of user property maps
// used by the RBAC user resolver to avoid repeated DB lookups per Check() call.
//
// When the cache reaches userCacheMaxSize entries, it is cleared before inserting
// the new entry. This simple eviction strategy avoids unbounded memory growth while
// keeping implementation complexity minimal.
type userCache struct {
	mu      sync.RWMutex
	entries map[uint64]map[string]interface{}
}

func newUserCache() userCache {
	return userCache{
		entries: make(map[uint64]map[string]interface{}),
	}
}

func (c *userCache) load(userID uint64) (map[string]interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.entries[userID]
	return v, ok
}

func (c *userCache) store(userID uint64, data map[string]interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.entries) >= userCacheMaxSize {
		c.entries = make(map[uint64]map[string]interface{})
	}
	c.entries[userID] = data
}

func (c *userCache) delete(userID uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, userID)
}

func (c *userCache) purge() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[uint64]map[string]interface{})
}
