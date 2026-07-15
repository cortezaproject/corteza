package tests

import (
	"context"
	"testing"
	"time"

	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/require"
)

func testProjectTasks(t *testing.T, s store.ProjectTasks) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func(projectID uint64) *types.ProjectTask {
			return &types.ProjectTask{
				ID:        id.Next(),
				ProjectID: projectID,
				Title:     "test task",
				CreatedAt: time.Now(),
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.ProjectTask) {
			req := require.New(t)
			req.NoError(s.TruncateProjectTasks(ctx))
			res := makeNew(id.Next())
			req.NoError(s.CreateProjectTask(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		req.NoError(s.TruncateProjectTasks(ctx))
		req.NoError(s.CreateProjectTask(ctx, makeNew(id.Next())))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, res := truncAndCreate(t)
		fetched, err := s.LookupProjectTaskByID(ctx, res.ID)
		req.NoError(err)
		req.Equal(res.ID, fetched.ID)
	})

	t.Run("update", func(t *testing.T) {
		req, res := truncAndCreate(t)
		res.Title = "updated"
		req.NoError(s.UpdateProjectTask(ctx, res))
		fetched, err := s.LookupProjectTaskByID(ctx, res.ID)
		req.NoError(err)
		req.Equal("updated", fetched.Title)
	})

	t.Run("delete", func(t *testing.T) {
		req, res := truncAndCreate(t)
		req.NoError(s.DeleteProjectTaskByID(ctx, res.ID))
	})
}
