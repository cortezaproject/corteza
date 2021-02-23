package tests

import (
	"context"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func testAuthConfirmedClients(t *testing.T, s store.AuthConfirmedClients) {
	var (
		ctx = context.Background()

		makeNew = func(c, u uint64) *types.AuthConfirmedClient {
			// minimum data set for new authConfirmedClient
			return &types.AuthConfirmedClient{
				ClientID:    c,
				UserID:      u,
				ConfirmedAt: time.Now(),
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.AuthConfirmedClient) {
			req := require.New(t)
			req.NoError(s.TruncateAuthConfirmedClients(ctx))
			res := makeNew(1, 2)
			req.NoError(s.CreateAuthConfirmedClient(ctx, res))
			return req, res
		}
	)

	t.Run("lookup by ID", func(t *testing.T) {
		req, acc := truncAndCreate(t)
		fetched, err := s.LookupAuthConfirmedClientByUserIDClientID(ctx, acc.UserID, acc.ClientID)
		req.NoError(err)
		req.NotNil(fetched)
	})

	t.Run("upsert", func(t *testing.T) {
		t.Run("existing", func(t *testing.T) {
			req, authConfirmedClient := truncAndCreate(t)
			req.NoError(s.UpsertAuthConfirmedClient(ctx, authConfirmedClient))
		})
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("by client and user ID", func(t *testing.T) {
			req, acc := truncAndCreate(t)
			req.NoError(s.DeleteAuthConfirmedClientByUserIDClientID(ctx, acc.UserID, acc.ClientID))
		})
	})

}
