package crs

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza-server/compose/crs/capabilities"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	ctx := context.Background()

	t.Run("invalid primary store", func(t *testing.T) {
		_, err := ComposeRecordStore(
			ctx,
			CRSConnectionWrap(0, "invalid://primary", capabilities.FullCapabilities()...),
			testDriver{},
		)
		require.Error(t, err)
	})

	t.Run("valid primary store", func(t *testing.T) {
		crs, err := ComposeRecordStore(
			ctx,
			CRSConnectionWrap(0, "noop://primary", capabilities.FullCapabilities()...),
			testDriver{},
		)
		require.NoError(t, err)

		require.Len(t, crs.stores, 0)
		require.NotNil(t, crs.primary)

	})
}

func TestStoreAdd(t *testing.T) {
	ctx := context.Background()

	t.Run("invalid external store", func(t *testing.T) {
		crs, err := ComposeRecordStore(
			ctx,
			CRSConnectionWrap(0, "noop://primary", capabilities.FullCapabilities()...),
			testDriver{},
		)
		require.NoError(t, err)

		sID := nextID()
		err = crs.AddStore(ctx, CRSConnectionWrap(sID, "invalid://external", capabilities.FullCapabilities()...))
		require.Error(t, err)

		require.NotNil(t, crs.primary)
		require.Len(t, crs.stores, 0)
	})

	t.Run("valid external store", func(t *testing.T) {
		crs, err := ComposeRecordStore(
			ctx,
			CRSConnectionWrap(0, "noop://primary", capabilities.FullCapabilities()...),
			testDriver{},
		)
		require.NoError(t, err)

		sID := nextID()
		err = crs.AddStore(ctx, CRSConnectionWrap(sID, "noop://external", capabilities.FullCapabilities()...))
		require.NoError(t, err)

		require.NotNil(t, crs.primary)
		require.Len(t, crs.stores, 1)
		require.NotNil(t, crs.stores[sID])
	})

	t.Run("duplicated external store", func(t *testing.T) {
		crs, err := ComposeRecordStore(
			ctx,
			CRSConnectionWrap(0, "noop://primary", capabilities.FullCapabilities()...),
			testDriver{},
		)
		require.NoError(t, err)

		sID := nextID()
		err = crs.AddStore(ctx, CRSConnectionWrap(sID, "noop://external", capabilities.FullCapabilities()...))
		require.NoError(t, err)

		err = crs.AddStore(ctx, CRSConnectionWrap(sID, "noop://external", capabilities.FullCapabilities()...))
		require.Error(t, err)

		require.NotNil(t, crs.primary)
		require.Len(t, crs.stores, 1)
		require.NotNil(t, crs.stores[sID])
	})
}

func TestStoreRemove(t *testing.T) {
	ctx := context.Background()

	t.Run("removing missing store", func(t *testing.T) {
		crs, err := ComposeRecordStore(
			ctx,
			CRSConnectionWrap(0, "noop://primary", capabilities.FullCapabilities()...),
			testDriver{},
		)
		require.NoError(t, err)

		sID := nextID()
		err = crs.AddStore(ctx, CRSConnectionWrap(sID, "noop://external", capabilities.FullCapabilities()...))
		require.NoError(t, err)

		err = crs.RemoveStore(ctx, 123)
		require.Error(t, err)
		require.Len(t, crs.stores, 1)
	})

	t.Run("removing store", func(t *testing.T) {
		crs, err := ComposeRecordStore(
			ctx,
			CRSConnectionWrap(0, "noop://primary", capabilities.FullCapabilities()...),
			testDriver{},
		)
		require.NoError(t, err)

		sID := nextID()
		err = crs.AddStore(ctx, CRSConnectionWrap(sID, "noop://external", capabilities.FullCapabilities()...))
		require.NoError(t, err)

		err = crs.RemoveStore(ctx, sID)
		require.NoError(t, err)
		require.Len(t, crs.stores, 0)
		require.NotNil(t, crs.primary)
	})
}
