package sass

import "sync"

type stylesheetCache struct {
	stylesheet map[string]string
	mu         sync.RWMutex
}

func newStylesheetCache() *stylesheetCache {
	return &stylesheetCache{
		stylesheet: make(map[string]string),
	}
}

func (c *stylesheetCache) Set(value map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.stylesheet = value
}

func (c *stylesheetCache) Get(key string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, exists := c.stylesheet[key]
	if !exists {
		return ""
	}
	return value
}

func (c *stylesheetCache) Empty() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.stylesheet) == 0
}
