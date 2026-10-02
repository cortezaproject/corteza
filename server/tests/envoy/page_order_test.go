package envoy

import (
	"context"
	"testing"
	"time"

	"github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/envoyx"
	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/store"
	"github.com/stretchr/testify/require"
)

// Pages must be exported in the order the page tree shows them (weight, then
// ID) so that importing the export recreates the same order.
func TestExportPageOrder(t *testing.T) {
	var (
		ctx = context.Background()
		req = require.New(t)
		now = time.Now()
	)

	cleanup(t)

	ns := &types.Namespace{ID: id.Next(), Slug: "page_order_ns", Name: "page order", CreatedAt: now}
	req.NoError(store.CreateComposeNamespace(ctx, defaultStore, ns))

	// Created in an order that matches neither weight nor the expected result;
	// pages with the same weight fall back to the ID order
	mk := func(handle string, weight int) *types.Page {
		return &types.Page{ID: id.Next(), NamespaceID: ns.ID, Handle: handle, Title: handle, Weight: weight, CreatedAt: now}
	}
	pages := []*types.Page{
		mk("third", 2),
		mk("first", 1),
		mk("fourth", 2),
		mk("second", 1),
	}
	// insert in reverse to make sure the physical order differs from the expected one
	for i := len(pages) - 1; i >= 0; i-- {
		req.NoError(store.CreateComposePage(ctx, defaultStore, pages[i]))
	}

	nodes, _, err := defaultEnvoy.Decode(ctx, envoyx.DecodeParams{
		Type: envoyx.DecodeTypeStore,
		Params: map[string]any{
			"storer": defaultStore,
			"dal":    defaultDal,
		},
		Filter: map[string]envoyx.ResourceFilter{
			types.PageResourceType: {
				Scope: envoyx.Scope{
					ResourceType: types.NamespaceResourceType,
					Identifiers:  envoyx.MakeIdentifiers(ns.Slug),
				},
			},
		},
	})
	req.NoError(err)

	var handles []string
	for _, n := range nodes {
		if n.ResourceType != types.PageResourceType {
			continue
		}
		handles = append(handles, n.Resource.(*types.Page).Handle)
	}

	req.Equal([]string{"first", "second", "third", "fourth"}, handles)
}
