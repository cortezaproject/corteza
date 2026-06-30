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

func testConnections(t *testing.T, s store.Connections) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func(handle string) *types.Connection {
			return &types.Connection{
				ID:        id.Next(),
				Handle:    handle,
				Status:    "active",
				CreatedAt: time.Now(),
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.Connection) {
			req := require.New(t)
			req.NoError(s.TruncateConnections(ctx))
			res := makeNew(string(rand.Bytes(10)))
			req.NoError(s.CreateConnection(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		req.NoError(s.TruncateConnections(ctx))
		req.NoError(s.CreateConnection(ctx, makeNew("TestConnection")))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, c := truncAndCreate(t)
		fetched, err := s.LookupConnectionByID(ctx, c.ID)
		req.NoError(err)
		req.Equal(c.ID, fetched.ID)
		req.Equal(c.Handle, fetched.Handle)
	})

	t.Run("update", func(t *testing.T) {
		req, c := truncAndCreate(t)
		c.Status = "disabled"
		req.NoError(s.UpdateConnection(ctx, c))
		fetched, err := s.LookupConnectionByID(ctx, c.ID)
		req.NoError(err)
		req.Equal("disabled", fetched.Status)
	})

	t.Run("delete", func(t *testing.T) {
		req, c := truncAndCreate(t)
		req.NoError(s.DeleteConnectionByID(ctx, c.ID))
	})
}
