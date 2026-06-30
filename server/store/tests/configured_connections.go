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

func testConfiguredConnections(t *testing.T, s store.ConfiguredConnections) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func(name string) *types.ConfiguredConnection {
			return &types.ConfiguredConnection{
				ID:        id.Next(),
				Name:      name,
				Status:    "active",
				CreatedAt: time.Now(),
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.ConfiguredConnection) {
			req := require.New(t)
			req.NoError(s.TruncateConfiguredConnections(ctx))
			res := makeNew(string(rand.Bytes(10)))
			req.NoError(s.CreateConfiguredConnection(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		req.NoError(s.TruncateConfiguredConnections(ctx))
		req.NoError(s.CreateConfiguredConnection(ctx, makeNew("TestConfiguredConnection")))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, c := truncAndCreate(t)
		fetched, err := s.LookupConfiguredConnectionByID(ctx, c.ID)
		req.NoError(err)
		req.Equal(c.ID, fetched.ID)
		req.Equal(c.Name, fetched.Name)
	})

	t.Run("update", func(t *testing.T) {
		req, c := truncAndCreate(t)
		c.Status = "disabled"
		req.NoError(s.UpdateConfiguredConnection(ctx, c))
		fetched, err := s.LookupConfiguredConnectionByID(ctx, c.ID)
		req.NoError(err)
		req.Equal("disabled", fetched.Status)
	})

	t.Run("delete", func(t *testing.T) {
		req, c := truncAndCreate(t)
		req.NoError(s.DeleteConfiguredConnectionByID(ctx, c.ID))
	})
}
