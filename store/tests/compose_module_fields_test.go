package tests

import (
	"context"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func testComposeModuleFields(t *testing.T, s composeModuleFieldsStore) {
	var (
		ctx = context.Background()
		req = require.New(t)

		moduleID = id.Next()

		makeNew = func(name, handle string) *types.ModuleField {
			// minimum data set for new composeModuleField
			return &types.ModuleField{
				ID:        id.Next(),
				ModuleID:  moduleID,
				CreatedAt: time.Now(),
				Name:      name,
			}
		}
	)

	t.Run("create", func(t *testing.T) {
		composeModuleField := makeNew("ComposeModuleFieldCRUD", "compose-moduleField-crud")
		req.NoError(s.CreateComposeModuleField(ctx, composeModuleField))
	})

	t.Run("create with duplicate handle", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("remove", func(t *testing.T) {
		composeModuleField := makeNew("remove", "remove")
		req.NoError(s.CreateComposeModuleField(ctx, composeModuleField))
		req.NoError(s.RemoveComposeModuleField(ctx))
	})

	t.Run("remove by ID", func(t *testing.T) {
		composeModuleField := makeNew("remove by id", "remove-by-id")
		req.NoError(s.CreateComposeModuleField(ctx, composeModuleField))
		req.NoError(s.RemoveComposeModuleField(ctx))
	})

	t.Run("update", func(t *testing.T) {
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
}
