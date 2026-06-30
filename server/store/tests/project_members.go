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

func testProjectMembers(t *testing.T, s store.ProjectMembers) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func(projectID, userID uint64) *types.ProjectMember {
			return &types.ProjectMember{
				ID:         id.Next(),
				ProjectID:  projectID,
				UserID:     userID,
				RolePreset: types.ProjectRoleMember,
				CreatedAt:  time.Now(),
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.ProjectMember) {
			req := require.New(t)
			req.NoError(s.TruncateProjectMembers(ctx))
			res := makeNew(id.Next(), id.Next())
			req.NoError(s.CreateProjectMember(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		req.NoError(s.TruncateProjectMembers(ctx))
		pm := makeNew(id.Next(), id.Next())
		req.NoError(s.CreateProjectMember(ctx, pm))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, pm := truncAndCreate(t)
		fetched, err := s.LookupProjectMemberByID(ctx, pm.ID)
		req.NoError(err)
		req.Equal(pm.ID, fetched.ID)
		req.Equal(pm.ProjectID, fetched.ProjectID)
	})

	t.Run("lookup by project and user", func(t *testing.T) {
		req, pm := truncAndCreate(t)
		fetched, err := s.LookupProjectMemberByProjectIDUserID(ctx, pm.ProjectID, pm.UserID)
		req.NoError(err)
		req.Equal(pm.ID, fetched.ID)
	})

	t.Run("update", func(t *testing.T) {
		req, pm := truncAndCreate(t)
		pm.RolePreset = types.ProjectRoleDeveloper
		req.NoError(s.UpdateProjectMember(ctx, pm))
		fetched, err := s.LookupProjectMemberByID(ctx, pm.ID)
		req.NoError(err)
		req.Equal(types.ProjectRoleDeveloper, fetched.RolePreset)
	})

	t.Run("delete", func(t *testing.T) {
		req, pm := truncAndCreate(t)
		req.NoError(s.DeleteProjectMemberByID(ctx, pm.ID))
	})
}
