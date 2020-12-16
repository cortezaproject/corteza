package store

import (
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/stretchr/testify/require"
)

func TestComposeNamespace_Merger(t *testing.T) {
	req := require.New(t)

	now := time.Time{}
	nowP := &time.Time{}

	empty := &types.Namespace{}
	full := &types.Namespace{
		Name:    "name",
		Slug:    "slug",
		Enabled: false,
		Meta: types.NamespaceMeta{
			Description: "dsc",
			Subtitle:    "sub",
		},

		CreatedAt: now,
		UpdatedAt: nowP,
		DeletedAt: nowP,
	}

	t.Run("merge on empty", func(t *testing.T) {
		c := mergeComposeNamespaces(empty, full)
		req.Equal("name", c.Name)
		req.Equal("slug", c.Slug)
		req.Equal("dsc", c.Meta.Description)
		req.Equal("sub", c.Meta.Subtitle)
		req.Equal(now, c.CreatedAt)
		req.Equal(nowP, c.UpdatedAt)
		req.Equal(nowP, c.DeletedAt)
	})

	t.Run("merge with empty", func(t *testing.T) {
		c := mergeComposeNamespaces(full, empty)
		req.Equal("name", c.Name)
		req.Equal("slug", c.Slug)
		req.Equal("dsc", c.Meta.Description)
		req.Equal("sub", c.Meta.Subtitle)
		req.Equal(now, c.CreatedAt)
		req.Equal(nowP, c.UpdatedAt)
		req.Equal(nowP, c.DeletedAt)
	})
}
