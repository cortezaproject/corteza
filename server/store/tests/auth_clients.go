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

func testAuthClients(t *testing.T, s store.AuthClients) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func(handle string) *types.AuthClient {
			// minimum data set for new authClient
			return &types.AuthClient{
				ID:        id.Next(),
				CreatedAt: time.Now(),
				Handle:    handle,
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.AuthClient) {
			req := require.New(t)
			req.NoError(s.TruncateAuthClients(ctx))
			res := makeNew(string(rand.Bytes(10)))
			req.NoError(s.CreateAuthClient(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		req.NoError(s.TruncateAuthClients(ctx))
		authClient := makeNew("AuthClientCRUD")
		req.NoError(s.CreateAuthClient(ctx, authClient))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, authClient := truncAndCreate(t)
		fetched, err := s.LookupAuthClientByID(ctx, authClient.ID)
		req.NoError(err)
		req.Equal(authClient.Handle, fetched.Handle)
		req.Equal(authClient.ID, fetched.ID)
		req.NotNil(fetched.CreatedAt)
		req.Nil(fetched.UpdatedAt)
		req.Nil(fetched.DeletedAt)
	})

	t.Run("update", func(t *testing.T) {
		req, authClient := truncAndCreate(t)
		authClient.Handle = "AuthClientCRUD-2"

		req.NoError(s.UpdateAuthClient(ctx, authClient))

		updated, err := s.LookupAuthClientByID(ctx, authClient.ID)
		req.NoError(err)
		req.Equal(authClient.Handle, updated.Handle)
	})

	t.Run("upsert", func(t *testing.T) {
		t.Run("existing", func(t *testing.T) {
			req, authClient := truncAndCreate(t)
			authClient.Handle = "AuthClientCRUD-2"

			req.NoError(s.UpsertAuthClient(ctx, authClient))

			updated, err := s.LookupAuthClientByID(ctx, authClient.ID)
			req.NoError(err)
			req.Equal(authClient.Handle, updated.Handle)
		})

		t.Run("new", func(t *testing.T) {
			authClient := makeNew("upsert me")
			authClient.Handle = "ComposeChartCRUD-2"

			req.NoError(s.UpsertAuthClient(ctx, authClient))

			upserted, err := s.LookupAuthClientByID(ctx, authClient.ID)
			req.NoError(err)
			req.Equal(authClient.Handle, upserted.Handle)
		})
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("by AuthClient", func(t *testing.T) {
			req, authClient := truncAndCreate(t)
			req.NoError(s.DeleteAuthClient(ctx, authClient))
			_, err := s.LookupAuthClientByID(ctx, authClient.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})

		t.Run("by ID", func(t *testing.T) {
			req, authClient := truncAndCreate(t)
			req.NoError(s.DeleteAuthClientByID(ctx, authClient.ID))
			_, err := s.LookupAuthClientByID(ctx, authClient.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})
	})

	t.Run("search", func(t *testing.T) {
		prefill := []*types.AuthClient{
			makeNew("/one-one"),
			makeNew("/one-two"),
			makeNew("/two-one"),
			makeNew("/two-two"),
			makeNew("/two-deleted"),
		}

		count := len(prefill)

		prefill[4].DeletedAt = &prefill[4].CreatedAt
		valid := count - 1

		req.NoError(s.TruncateAuthClients(ctx))
		req.NoError(s.CreateAuthClient(ctx, prefill...))

		// search for all valid
		set, f, err := s.SearchAuthClients(ctx, types.AuthClientFilter{})
		req.NoError(err)
		req.Len(set, valid) // we've deleted one

		// search for ALL
		set, f, err = s.SearchAuthClients(ctx, types.AuthClientFilter{Deleted: filter.StateInclusive})
		req.NoError(err)
		req.Len(set, count) // we've deleted one

		// search for deleted only
		set, f, err = s.SearchAuthClients(ctx, types.AuthClientFilter{Deleted: filter.StateExclusive})
		req.NoError(err)
		req.Len(set, 1) // we've deleted one

		set, f, err = s.SearchAuthClients(ctx, types.AuthClientFilter{Handle: "/two-one"})
		req.NoError(err)
		req.Len(set, 1)

		_ = f // dummy
	})
}
