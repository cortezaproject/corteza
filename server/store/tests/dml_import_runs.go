package tests

import (
	"context"
	"testing"

	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/require"
)

func testDmlImportRuns(t *testing.T, s store.DmlImportRuns) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func() *types.DmlImportRun {
			return &types.DmlImportRun{
				ID:           id.Next(),
				ConnectionID: 42,
				MappingID:    1,
				Method:       "full",
				Status:       "pending",
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.DmlImportRun) {
			req := require.New(t)
			req.NoError(s.TruncateDmlImportRuns(ctx))
			res := makeNew()
			req.NoError(s.CreateDmlImportRun(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		req.NoError(s.TruncateDmlImportRuns(ctx))
		req.NoError(s.CreateDmlImportRun(ctx, makeNew()))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, r := truncAndCreate(t)
		fetched, err := s.LookupDmlImportRunByID(ctx, r.ID)
		req.NoError(err)
		req.Equal(r.ID, fetched.ID)
		req.Equal(r.Status, fetched.Status)
	})

	t.Run("update", func(t *testing.T) {
		req, r := truncAndCreate(t)
		r.Status = "done"
		req.NoError(s.UpdateDmlImportRun(ctx, r))
		fetched, err := s.LookupDmlImportRunByID(ctx, r.ID)
		req.NoError(err)
		req.Equal("done", fetched.Status)
	})

	t.Run("delete", func(t *testing.T) {
		req, r := truncAndCreate(t)
		req.NoError(s.DeleteDmlImportRunByID(ctx, r.ID))
	})
}
