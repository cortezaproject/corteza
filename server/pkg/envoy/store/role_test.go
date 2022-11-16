package store

import (
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/stretchr/testify/require"
)

func TestRole_Merger(t *testing.T) {
	req := require.New(t)

	now := time.Time{}
	nowP := &time.Time{}

	empty := &types.Role{}
	full := &types.Role{
		Name:       "name",
		Handle:     "handle",
		CreatedAt:  now,
		UpdatedAt:  nowP,
		ArchivedAt: nowP,
		DeletedAt:  nowP,
	}

	t.Run("merge on empty", func(t *testing.T) {
		c := mergeRoles(empty, full)
		req.Equal("name", c.Name)
		req.Equal("handle", c.Handle)
		req.Equal(now, c.CreatedAt)
		req.Equal(nowP, c.UpdatedAt)
		req.Equal(nowP, c.ArchivedAt)
		req.Equal(nowP, c.DeletedAt)
	})

	t.Run("merge with empty", func(t *testing.T) {
		c := mergeRoles(full, empty)
		req.Equal("name", c.Name)
		req.Equal("handle", c.Handle)
		req.Equal(now, c.CreatedAt)
		req.Equal(nowP, c.UpdatedAt)
		req.Equal(nowP, c.ArchivedAt)
		req.Equal(nowP, c.DeletedAt)
	})
}
