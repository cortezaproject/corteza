package crs

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/stretchr/testify/require"
)

var (
	defaultExternalCRS uint64 = nextID()
)

func TestLoadModules(t *testing.T) {
	ctx := context.Background()

	t.Run("no supporting CRS", func(t *testing.T) {
		crs := initCRS(ctx)

		err := crs.ReloadModules(ctx, &types.Module{
			ID:    nextID(),
			Store: types.CRSDef{ComposeRecordStoreID: 999},
		})
		require.Error(t, err)
	})

	t.Run("module added", func(t *testing.T) {
		crs := initCRS(ctx)

		mod := defaultModule()

		err := crs.ReloadModules(ctx, mod)
		require.NoError(t, err)

		model := crs.lookupModel(mod)
		require.NotNil(t, model)

		// 8 sys, 1 custom
		require.Len(t, model.Attributes, 8+1)

		require.Equal(t, defaultExternalCRS, model.StoreID)
	})

	t.Run("duplicate module added", func(t *testing.T) {
		crs := initCRS(ctx)

		mod := defaultModule()

		err := crs.ReloadModules(ctx, mod, mod)
		require.Error(t, err)
	})
}

func TestRemoveModules(t *testing.T) {
	ctx := context.Background()

	t.Run("module not found", func(t *testing.T) {
		crs := initCRS(ctx)

		err := crs.RemoveModules(ctx, defaultModule())
		require.Error(t, err)
	})

	t.Run("module removed", func(t *testing.T) {
		crs := initCRS(ctx)
		mod := defaultModule()

		err := crs.ReloadModules(ctx, mod)
		require.NoError(t, err)

		err = crs.RemoveModules(ctx, mod)
		require.NoError(t, err)

		require.Nil(t, crs.lookupModel(mod))
	})
}
