package tests

import (
	"context"
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/pkg/rand"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/stretchr/testify/require"
)

func testDalConnections(t *testing.T, s store.DalConnections) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func(handle string) *types.DalConnection {
			// minimum data set for new dalConnection
			return &types.DalConnection{
				ID:        id.Next(),
				CreatedAt: time.Now(),
				Handle:    handle,
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.DalConnection) {
			req := require.New(t)
			req.NoError(s.TruncateDalConnections(ctx))
			res := makeNew(string(rand.Bytes(10)))
			req.NoError(s.CreateDalConnection(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		req.NoError(s.TruncateDalConnections(ctx))
		dalConnection := makeNew("DalConnectionCRUD")
		req.NoError(s.CreateDalConnection(ctx, dalConnection))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, dalConnection := truncAndCreate(t)
		fetched, err := s.LookupDalConnectionByID(ctx, dalConnection.ID)
		req.NoError(err)
		req.Equal(dalConnection.Handle, fetched.Handle)
		req.Equal(dalConnection.ID, fetched.ID)
		req.NotNil(fetched.CreatedAt)
		req.Nil(fetched.UpdatedAt)
		req.Nil(fetched.DeletedAt)
	})

	t.Run("update", func(t *testing.T) {
		req, dalConnection := truncAndCreate(t)
		dalConnection.Handle = "DalConnectionCRUD-2"

		req.NoError(s.UpdateDalConnection(ctx, dalConnection))

		updated, err := s.LookupDalConnectionByID(ctx, dalConnection.ID)
		req.NoError(err)
		req.Equal(dalConnection.Handle, updated.Handle)
	})

	t.Run("upsert", func(t *testing.T) {
		t.Run("existing", func(t *testing.T) {
			req, dalConnection := truncAndCreate(t)
			dalConnection.Handle = "DalConnectionCRUD-2"

			req.NoError(s.UpsertDalConnection(ctx, dalConnection))

			updated, err := s.LookupDalConnectionByID(ctx, dalConnection.ID)
			req.NoError(err)
			req.Equal(dalConnection.Handle, updated.Handle)
		})

		t.Run("new", func(t *testing.T) {
			dalConnection := makeNew("upsert me")
			dalConnection.Handle = "ComposeChartCRUD-2"

			req.NoError(s.UpsertDalConnection(ctx, dalConnection))

			upserted, err := s.LookupDalConnectionByID(ctx, dalConnection.ID)
			req.NoError(err)
			req.Equal(dalConnection.Handle, upserted.Handle)
		})
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("by DalConnection", func(t *testing.T) {
			req, dalConnection := truncAndCreate(t)
			req.NoError(s.DeleteDalConnection(ctx, dalConnection))
			_, err := s.LookupDalConnectionByID(ctx, dalConnection.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})

		t.Run("by ID", func(t *testing.T) {
			req, dalConnection := truncAndCreate(t)
			req.NoError(s.DeleteDalConnectionByID(ctx, dalConnection.ID))
			_, err := s.LookupDalConnectionByID(ctx, dalConnection.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})
	})

	t.Run("search", func(t *testing.T) {
		prefill := []*types.DalConnection{
			makeNew("/one-one"),
			makeNew("/one-two"),
			makeNew("/two-one"),
			makeNew("/two-two"),
			makeNew("/two-deleted"),
		}

		count := len(prefill)

		prefill[4].DeletedAt = &prefill[4].CreatedAt
		valid := count - 1

		req.NoError(s.TruncateDalConnections(ctx))
		req.NoError(s.CreateDalConnection(ctx, prefill...))

		// search for all valid
		set, _, err := s.SearchDalConnections(ctx, types.DalConnectionFilter{})
		req.NoError(err)
		req.Len(set, valid) // we've deleted one

		// search for ALL
		set, _, err = s.SearchDalConnections(ctx, types.DalConnectionFilter{Deleted: filter.StateInclusive})
		req.NoError(err)
		req.Len(set, count) // we've deleted one

		// search for deleted only
		set, _, err = s.SearchDalConnections(ctx, types.DalConnectionFilter{Deleted: filter.StateExclusive})
		req.NoError(err)
		req.Len(set, 1) // we've deleted one

		set, _, err = s.SearchDalConnections(ctx, types.DalConnectionFilter{Handle: "/two-one"})
		req.NoError(err)
		req.Len(set, 1)
	})
}
