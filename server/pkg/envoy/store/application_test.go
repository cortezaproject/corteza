package store

import (
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/stretchr/testify/require"
)

func TestApplication_Merger(t *testing.T) {
	req := require.New(t)

	ub := &types.ApplicationUnify{
		Config: "config",
		Icon:   "icon",
	}

	now := time.Time{}
	nowP := &time.Time{}

	empty := &types.Application{}
	full := &types.Application{
		Name:      "name",
		OwnerID:   1,
		Enabled:   true,
		Unify:     ub,
		CreatedAt: now,
		UpdatedAt: nowP,
		DeletedAt: nowP,
	}

	t.Run("merge on empty", func(t *testing.T) {
		c := mergeApplications(empty, full)
		req.Equal("name", c.Name)
		req.Equal(uint64(1), c.OwnerID)
		req.False(c.Enabled)
		req.Equal(ub, c.Unify)
		req.Equal(now, c.CreatedAt)
		req.Equal(nowP, c.UpdatedAt)
		req.Equal(nowP, c.DeletedAt)
	})

	t.Run("merge with empty", func(t *testing.T) {
		c := mergeApplications(full, empty)
		req.Equal("name", c.Name)
		req.Equal(uint64(0), c.OwnerID)
		req.True(c.Enabled)
		req.Equal(ub, c.Unify)
		req.Equal(now, c.CreatedAt)
		req.Equal(nowP, c.UpdatedAt)
		req.Equal(nowP, c.DeletedAt)
	})
}
