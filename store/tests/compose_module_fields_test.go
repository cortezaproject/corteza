package tests

import (
	"context"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/rand"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func testComposeModuleFields(t *testing.T, s store.ComposeModuleFields) {
	var (
		ctx = context.Background()

		moduleID = id.Next()

		makeNew = func(name, label string) *types.ModuleField {
			// minimum data set for new composeModuleField
			return &types.ModuleField{
				ID:        id.Next(),
				ModuleID:  moduleID,
				CreatedAt: time.Now(),
				Name:      name,
				Label:     label,
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.ModuleField) {
			req := require.New(t)
			req.NoError(s.TruncateComposeModuleFields(ctx))
			res := makeNew(string(rand.Bytes(10)), string(rand.Bytes(10)))
			req.NoError(s.CreateComposeModuleField(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.TruncateComposeModuleFields(ctx))
		composeModuleField := makeNew("ComposeModuleFieldCRUD", "compose-moduleField-crud")
		req.NoError(s.CreateComposeModuleField(ctx, composeModuleField))
	})

	t.Run("create with duplicate name", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.TruncateComposeModuleFields(ctx))
		t.Skip("not implemented")
	})

	t.Run("update", func(t *testing.T) {
		req := require.New(t)
		composeModuleField := makeNew("update me", "update-me")
		req.NoError(s.CreateComposeModuleField(ctx, composeModuleField))

		composeModuleField = &types.ModuleField{
			ID:        composeModuleField.ID,
			CreatedAt: composeModuleField.CreatedAt,
			Name:      "ComposeModuleFieldCRUD+2",
		}
		req.NoError(s.UpdateComposeModuleField(ctx, composeModuleField))
	})

	t.Run("update with duplicate handle", func(t *testing.T) {
		t.Skip("not implemented")
	})


	t.Run("delete", func(t *testing.T) {
		t.Run("by Field", func(t *testing.T) {
			req, fld := truncAndCreate(t)
			req.NoError(s.DeleteComposeModuleField(ctx, fld))
			fetched, _, err := s.SearchComposeModuleFields(ctx, types.ModuleFieldFilter{ModuleID: []uint64{fld.ModuleID}})
			req.NoError(err)
			req.Empty(fetched)
		})

		t.Run("by ID", func(t *testing.T) {
			req, fld := truncAndCreate(t)
			req.NoError(s.DeleteComposeModuleFieldByID(ctx, fld.ID))
			fetched, _, err := s.SearchComposeModuleFields(ctx, types.ModuleFieldFilter{ModuleID: []uint64{fld.ModuleID}})
			req.NoError(err)
			req.Empty(fetched)
		})
	})
}
