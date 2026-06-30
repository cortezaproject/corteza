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

func testDmlMappings(t *testing.T, s store.DmlMappings) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func(sourceIdent string) *types.DmlMapping {
			return &types.DmlMapping{
				ID:           id.Next(),
				ConnectionID: 42,
				SourceIdent:  sourceIdent,
				ModuleHandle: string(rand.Bytes(8)),
				CreatedAt:    time.Now(),
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.DmlMapping) {
			req := require.New(t)
			req.NoError(s.TruncateDmlMappings(ctx))
			res := makeNew(string(rand.Bytes(10)))
			req.NoError(s.CreateDmlMapping(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		req.NoError(s.TruncateDmlMappings(ctx))
		req.NoError(s.CreateDmlMapping(ctx, makeNew("test_source")))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, m := truncAndCreate(t)
		fetched, err := s.LookupDmlMappingByID(ctx, m.ID)
		req.NoError(err)
		req.Equal(m.ID, fetched.ID)
		req.Equal(m.SourceIdent, fetched.SourceIdent)
	})

	t.Run("update", func(t *testing.T) {
		req, m := truncAndCreate(t)
		m.ModuleHandle = "updated-module"
		req.NoError(s.UpdateDmlMapping(ctx, m))
		fetched, err := s.LookupDmlMappingByID(ctx, m.ID)
		req.NoError(err)
		req.Equal("updated-module", fetched.ModuleHandle)
	})

	t.Run("delete", func(t *testing.T) {
		req, m := truncAndCreate(t)
		req.NoError(s.DeleteDmlMappingByID(ctx, m.ID))
	})
}
