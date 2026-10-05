package service

import (
	"context"
	"testing"

	"github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/store"
	"github.com/stretchr/testify/require"
)

// moduleStoreStub knows modules of namespace 1 by handle.
type moduleStoreStub struct {
	store.ComposeModules
	ids map[string]uint64
}

func (s moduleStoreStub) LookupComposeModuleByNamespaceIDHandle(_ context.Context, namespaceID uint64, handle string) (*types.Module, error) {
	if id, ok := s.ids[handle]; ok && namespaceID == 1 {
		return &types.Module{ID: id, NamespaceID: 1, Handle: handle}, nil
	}
	return nil, store.ErrNotFound
}

func TestCheckCustomBlocks(t *testing.T) {
	var (
		ctx     = context.Background()
		modules = moduleStoreStub{ids: map[string]uint64{"task": 42}}
		check   = func(bb types.PageBlocks) error { return checkCustomBlocks(ctx, modules, 1, bb) }
		custom  = func(opts map[string]any) types.PageBlocks {
			return types.PageBlocks{{Kind: "Content"}, {Kind: "Custom", Options: opts}}
		}
	)

	require.NoError(t, check(custom(map[string]any{"applicationID": "123"})))
	require.NoError(t, check(custom(map[string]any{"source": "<p>hi</p>"})))

	require.ErrorContains(t, check(custom(map[string]any{})), "block 2 (Custom): the source is empty")
	require.ErrorContains(t, check(custom(map[string]any{"source": "<script>fetch('/x')</script>"})), "connect-src 'none'")

	// HTML of its own is checked even beside an application.
	require.ErrorContains(t, check(custom(map[string]any{"applicationID": "123", "source": "<script>alert(1)</script>"})), "alert()")

	// A script from an origin the block lists is admitted, and the list is
	// stored as bare origins.
	blocks := custom(map[string]any{
		"source":  `<script src="https://cdn.jsdelivr.net/npm/chart.js@4.4.0"></script>`,
		"origins": []any{"https://CDN.jsdelivr.net/", "https://cdn.jsdelivr.net"},
	})
	require.NoError(t, check(blocks))
	require.Equal(t, []string{"https://cdn.jsdelivr.net"}, blocks[1].Options["origins"])

	require.ErrorContains(t, check(custom(map[string]any{"source": `<script src="https://cdn.jsdelivr.net/npm/x.js"></script>`})), "allowed origins")
	require.ErrorContains(t, check(custom(map[string]any{"source": "<p>hi</p>", "origins": []any{"http://evil.test"}})), "not https")
}

func TestCheckCustomBlocksResolvesWhatItDeclares(t *testing.T) {
	ctx := context.Background()
	modules := moduleStoreStub{ids: map[string]uint64{"task": 42}}
	block := func(opts map[string]any) types.PageBlocks {
		opts["source"] = "<p>hi</p>"
		return types.PageBlocks{{Kind: "Custom", Options: opts}}
	}

	// Each module is stored beside its ID, whatever the caller sent there.
	bb := block(map[string]any{"modules": []any{"task"}, "moduleIDs": map[string]any{"task": "999"}})
	require.NoError(t, checkCustomBlocks(ctx, modules, 1, bb))
	require.Equal(t, map[string]string{"task": "42"}, bb[0].Options["moduleIDs"])

	err := checkCustomBlocks(ctx, modules, 1, block(map[string]any{"modules": []any{"lead"}}))
	require.ErrorContains(t, err, `there is no module "lead" in this page's namespace`)

	// A module of another namespace is no module of this page's.
	err = checkCustomBlocks(ctx, modules, 2, block(map[string]any{"modules": []any{"task"}}))
	require.ErrorContains(t, err, `no module "task"`)

	err = checkCustomBlocks(ctx, modules, 1, block(map[string]any{"modules": []any{"task"}, "writes": []any{"lead"}}))
	require.ErrorContains(t, err, `may change module "lead" but does not read it`)

	err = checkCustomBlocks(ctx, modules, 1, block(map[string]any{"modules": []any{"task"}, "deletes": []any{"lead"}}))
	require.ErrorContains(t, err, `may delete from module "lead" but does not read it`)
}
