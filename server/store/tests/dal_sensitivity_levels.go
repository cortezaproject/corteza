package tests

import (
	"context"
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/rand"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/stretchr/testify/require"
)

func testDalSensitivityLevels(t *testing.T, s store.DalSensitivityLevels) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func(handle string) *types.DalSensitivityLevel {
			// minimum data set for new dalSensitivityLevel
			return &types.DalSensitivityLevel{
				ID:        id.Next(),
				CreatedAt: time.Now(),
				Handle:    handle,
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.DalSensitivityLevel) {
			req := require.New(t)
			req.NoError(s.TruncateDalSensitivityLevels(ctx))
			res := makeNew(string(rand.Bytes(10)))
			req.NoError(s.CreateDalSensitivityLevel(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		dalSensitivityLevel := makeNew("DalSensitivityLevelCRUD")
		req.NoError(s.CreateDalSensitivityLevel(ctx, dalSensitivityLevel))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, dalSensitivityLevel := truncAndCreate(t)
		fetched, err := s.LookupDalSensitivityLevelByID(ctx, dalSensitivityLevel.ID)
		req.NoError(err)
		req.Equal(dalSensitivityLevel.Handle, fetched.Handle)
		req.Equal(dalSensitivityLevel.ID, fetched.ID)
		req.NotNil(fetched.CreatedAt)
		req.Nil(fetched.UpdatedAt)
		req.Nil(fetched.DeletedAt)
	})

	t.Run("update", func(t *testing.T) {
		req, dalSensitivityLevel := truncAndCreate(t)
		dalSensitivityLevel.Handle = "DalSensitivityLevelCRUD-2"

		req.NoError(s.UpdateDalSensitivityLevel(ctx, dalSensitivityLevel))

		updated, err := s.LookupDalSensitivityLevelByID(ctx, dalSensitivityLevel.ID)
		req.NoError(err)
		req.Equal(dalSensitivityLevel.Handle, updated.Handle)
	})

	t.Run("upsert", func(t *testing.T) {
		t.Run("existing", func(t *testing.T) {
			req, dalSensitivityLevel := truncAndCreate(t)
			dalSensitivityLevel.Handle = "DalSensitivityLevelCRUD-2"

			req.NoError(s.UpsertDalSensitivityLevel(ctx, dalSensitivityLevel))

			updated, err := s.LookupDalSensitivityLevelByID(ctx, dalSensitivityLevel.ID)
			req.NoError(err)
			req.Equal(dalSensitivityLevel.Handle, updated.Handle)
		})

		t.Run("new", func(t *testing.T) {
			dalSensitivityLevel := makeNew("upsert me")
			dalSensitivityLevel.Handle = "ComposeChartCRUD-2"

			req.NoError(s.UpsertDalSensitivityLevel(ctx, dalSensitivityLevel))

			upserted, err := s.LookupDalSensitivityLevelByID(ctx, dalSensitivityLevel.ID)
			req.NoError(err)
			req.Equal(dalSensitivityLevel.Handle, upserted.Handle)
		})
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("by DalSensitivityLevel", func(t *testing.T) {
			req, dalSensitivityLevel := truncAndCreate(t)
			req.NoError(s.DeleteDalSensitivityLevel(ctx, dalSensitivityLevel))
			_, err := s.LookupDalSensitivityLevelByID(ctx, dalSensitivityLevel.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})

		t.Run("by ID", func(t *testing.T) {
			req, dalSensitivityLevel := truncAndCreate(t)
			req.NoError(s.DeleteDalSensitivityLevelByID(ctx, dalSensitivityLevel.ID))
			_, err := s.LookupDalSensitivityLevelByID(ctx, dalSensitivityLevel.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})
	})

	t.Run("search", func(t *testing.T) {
		prefill := []*types.DalSensitivityLevel{
			makeNew("/one-one"),
			makeNew("/one-two"),
			makeNew("/two-one"),
			makeNew("/two-two"),
			makeNew("/two-deleted"),
		}

		count := len(prefill)

		prefill[4].DeletedAt = &prefill[4].CreatedAt
		valid := count - 1

		req.NoError(s.TruncateDalSensitivityLevels(ctx))
		req.NoError(s.CreateDalSensitivityLevel(ctx, prefill...))

		// search for all valid
		set, _, err := s.SearchDalSensitivityLevels(ctx, types.DalSensitivityLevelFilter{})
		req.NoError(err)
		req.Len(set, valid) // we've deleted one

		// search for ALL
		set, _, err = s.SearchDalSensitivityLevels(ctx, types.DalSensitivityLevelFilter{Deleted: filter.StateInclusive})
		req.NoError(err)
		req.Len(set, count) // we've deleted one

		// search for deleted only
		set, _, err = s.SearchDalSensitivityLevels(ctx, types.DalSensitivityLevelFilter{Deleted: filter.StateExclusive})
		req.NoError(err)
		req.Len(set, 1) // we've deleted one
	})
}
