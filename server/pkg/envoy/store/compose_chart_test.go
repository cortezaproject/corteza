package store

import (
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/stretchr/testify/require"
)

func TestComposeChart_Merger(t *testing.T) {
	req := require.New(t)

	now := time.Time{}
	nowP := &time.Time{}

	empty := &types.Chart{}
	full := &types.Chart{
		Handle: "handle",
		Name:   "name",
		Config: types.ChartConfig{
			Reports:     []*types.ChartConfigReport{&types.ChartConfigReport{}},
			ColorScheme: "colorScheme",
		},
		NamespaceID: 1,
		CreatedAt:   now,
		UpdatedAt:   nowP,
		DeletedAt:   nowP,
	}

	t.Run("merge on empty", func(t *testing.T) {
		c := mergeComposeChart(empty, full)
		req.Equal("name", c.Name)
		req.Equal("handle", c.Handle)
		req.Len(c.Config.Reports, 1)
		req.Equal("colorScheme", c.Config.ColorScheme)
		req.Equal(uint64(1), c.NamespaceID)
		req.Equal(now, c.CreatedAt)
		req.Equal(nowP, c.UpdatedAt)
		req.Equal(nowP, c.DeletedAt)
	})

	t.Run("merge with empty", func(t *testing.T) {
		c := mergeComposeChart(full, empty)
		req.Equal("name", c.Name)
		req.Equal("handle", c.Handle)
		req.Len(c.Config.Reports, 1)
		req.Equal("colorScheme", c.Config.ColorScheme)
		req.Equal(uint64(0), c.NamespaceID)
		req.Equal(now, c.CreatedAt)
		req.Equal(nowP, c.UpdatedAt)
		req.Equal(nowP, c.DeletedAt)
	})
}
