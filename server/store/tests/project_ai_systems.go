package tests

import (
	"context"
	"testing"
	"time"

	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/pkg/rand"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/require"
)

func testProjectGroups(t *testing.T, s store.ProjectGroups) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func(handle string) *types.ProjectGroup {
			return &types.ProjectGroup{
				ID:        id.Next(),
				ProjectID: id.Next(),
				CreatedAt: time.Now(),
				Handle:    handle,
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.ProjectGroup) {
			req := require.New(t)
			req.NoError(s.TruncateProjectGroups(ctx))
			res := makeNew(string(rand.Bytes(10)))
			req.NoError(s.CreateProjectGroup(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		req.NoError(s.TruncateProjectGroups(ctx))
		pg := makeNew("ProjectGroupCRUD")
		req.NoError(s.CreateProjectGroup(ctx, pg))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, pg := truncAndCreate(t)
		fetched, err := s.LookupProjectGroupByID(ctx, pg.ID)
		req.NoError(err)
		req.Equal(pg.Handle, fetched.Handle)
		req.Equal(pg.ID, fetched.ID)
	})

	t.Run("update", func(t *testing.T) {
		req, pg := truncAndCreate(t)
		pg.Handle = "updated-handle"
		req.NoError(s.UpdateProjectGroup(ctx, pg))
		fetched, err := s.LookupProjectGroupByID(ctx, pg.ID)
		req.NoError(err)
		req.Equal("updated-handle", fetched.Handle)
	})

	t.Run("delete", func(t *testing.T) {
		req, pg := truncAndCreate(t)
		req.NoError(s.DeleteProjectGroupByID(ctx, pg.ID))
	})
}
