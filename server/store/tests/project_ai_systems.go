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

func testProjectAiSystems(t *testing.T, s store.ProjectAiSystems) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func(handle string) *types.ProjectAiSystem {
			return &types.ProjectAiSystem{
				ID:        id.Next(),
				ProjectID: id.Next(),
				CreatedAt: time.Now(),
				Handle:    handle,
				RiskClass: "high",
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.ProjectAiSystem) {
			req := require.New(t)
			req.NoError(s.TruncateProjectAiSystems(ctx))
			res := makeNew(string(rand.Bytes(10)))
			req.NoError(s.CreateProjectAiSystem(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		req.NoError(s.TruncateProjectAiSystems(ctx))
		pg := makeNew("ProjectAiSystemCRUD")
		req.NoError(s.CreateProjectAiSystem(ctx, pg))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, pg := truncAndCreate(t)
		fetched, err := s.LookupProjectAiSystemByID(ctx, pg.ID)
		req.NoError(err)
		req.Equal(pg.Handle, fetched.Handle)
		req.Equal(pg.ID, fetched.ID)
		req.Equal(pg.RiskClass, fetched.RiskClass)
	})

	t.Run("update", func(t *testing.T) {
		req, pg := truncAndCreate(t)
		pg.Handle = "updated-handle"
		pg.RiskClass = "limited"
		req.NoError(s.UpdateProjectAiSystem(ctx, pg))
		fetched, err := s.LookupProjectAiSystemByID(ctx, pg.ID)
		req.NoError(err)
		req.Equal("updated-handle", fetched.Handle)
		req.Equal("limited", fetched.RiskClass)
	})

	t.Run("delete", func(t *testing.T) {
		req, pg := truncAndCreate(t)
		req.NoError(s.DeleteProjectAiSystemByID(ctx, pg.ID))
	})
}
