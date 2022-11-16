package corredor

import (
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	userLookupCache struct {
		err  error
		user *types.User
	}

	userLookupCacheMap map[string]userLookupCache
)

func (m userLookupCacheMap) lookup(key string, lookup func() (*types.User, error)) (*types.User, error) {
	if c, ok := m[key]; ok {
		return c.user, c.err
	}

	c := userLookupCache{}
	c.user, c.err = lookup()

	m[key] = c

	return c.user, c.err
}
