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

func testProjects(t *testing.T, s store.Projects) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func(handle string) *types.Project {
			return &types.Project{
				ID:        id.Next(),
				CreatedAt: time.Now(),
				Handle:    handle,
				Status:    types.ProjectStatusActive,
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.Project) {
			req := require.New(t)
			req.NoError(s.TruncateProjects(ctx))
			res := makeNew(string(rand.Bytes(10)))
			req.NoError(s.CreateProject(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		req.NoError(s.TruncateProjects(ctx))
		project := makeNew("ProjectCRUD")
		req.NoError(s.CreateProject(ctx, project))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, project := truncAndCreate(t)
		fetched, err := s.LookupProjectByID(ctx, project.ID)
		req.NoError(err)
		req.Equal(project.Handle, fetched.Handle)
		req.Equal(project.ID, fetched.ID)
	})

	t.Run("update", func(t *testing.T) {
		req, project := truncAndCreate(t)
		project.Status = types.ProjectStatusArchived
		req.NoError(s.UpdateProject(ctx, project))
		fetched, err := s.LookupProjectByID(ctx, project.ID)
		req.NoError(err)
		req.Equal(types.ProjectStatusArchived, fetched.Status)
	})

	t.Run("delete", func(t *testing.T) {
		req, project := truncAndCreate(t)
		req.NoError(s.DeleteProjectByID(ctx, project.ID))
	})
}
