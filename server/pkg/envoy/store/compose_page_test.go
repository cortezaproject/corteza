package store

import (
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/stretchr/testify/require"
)

func TestComposePage_Merger(t *testing.T) {
	req := require.New(t)

	now := time.Time{}
	nowP := &time.Time{}

	empty := &types.Page{}
	full := &types.Page{
		SelfID:      1,
		NamespaceID: 2,
		ModuleID:    3,
		Handle:      "handle",
		Title:       "title",
		Description: "description",
		Blocks:      types.PageBlocks{types.PageBlock{}},
		Children:    types.PageSet{&types.Page{}},
		Weight:      4,

		CreatedAt: now,
		UpdatedAt: nowP,
		DeletedAt: nowP,
	}

	t.Run("merge on empty", func(t *testing.T) {
		c := mergeComposePage(empty, full)
		req.Equal(uint64(1), c.SelfID)
		req.Equal(uint64(2), c.NamespaceID)
		req.Equal(uint64(3), c.ModuleID)

		req.Equal("handle", c.Handle)
		req.Equal("title", c.Title)
		req.Equal("description", c.Description)
		req.Len(c.Blocks, 1)
		req.Len(c.Children, 1)
		req.Equal(4, c.Weight)

		req.Equal(now, c.CreatedAt)
		req.Equal(nowP, c.UpdatedAt)
		req.Equal(nowP, c.DeletedAt)
	})

	t.Run("merge with empty", func(t *testing.T) {
		c := mergeComposePage(full, empty)
		req.Equal(uint64(0), c.SelfID)
		req.Equal(uint64(0), c.NamespaceID)
		req.Equal(uint64(0), c.ModuleID)

		req.Equal("handle", c.Handle)
		req.Equal("title", c.Title)
		req.Equal("description", c.Description)
		req.Len(c.Blocks, 1)
		req.Len(c.Children, 1)
		req.Equal(0, c.Weight)

		req.Equal(now, c.CreatedAt)
		req.Equal(nowP, c.UpdatedAt)
		req.Equal(nowP, c.DeletedAt)
	})
}
