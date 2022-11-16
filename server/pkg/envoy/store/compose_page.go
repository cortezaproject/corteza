package store

import (
	"context"

	atypes "github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/store"
)

type (
	composePage struct {
		cfg *EncoderConfig

		res *resource.ComposePage
		pg  *types.Page

		relNS     *types.Namespace
		relMod    *types.Module
		relParent *types.Page

		relWfs    map[string]*atypes.Workflow
		relMods   map[string]*types.Module
		relCharts map[string]*types.Chart
	}
)

// mergeComposePage merges b into a, prioritising a
func mergeComposePage(a, b *types.Page) *types.Page {
	c := a.Clone()

	c.SelfID = b.SelfID
	c.NamespaceID = b.NamespaceID
	c.ModuleID = b.ModuleID
	c.Weight = b.Weight
	if c.Handle == "" {
		c.Handle = b.Handle
	}
	if c.Title == "" {
		c.Title = b.Title
	}
	if c.Description == "" {
		c.Description = b.Description
	}
	if len(c.Blocks) <= 0 {
		c.Blocks = b.Blocks
	}
	if len(c.Children) <= 0 {
		c.Children = b.Children
	}

	if c.CreatedAt.IsZero() {
		c.CreatedAt = b.CreatedAt
	}
	if c.UpdatedAt == nil {
		c.UpdatedAt = b.UpdatedAt
	}
	if c.DeletedAt == nil {
		c.DeletedAt = b.DeletedAt
	}

	return c
}

// findComposePage looks for the page in the resources & the store
//
// Provided resources are prioritized.
func findComposePage(ctx context.Context, s store.Storer, nsID uint64, rr resource.InterfaceSet, ii resource.Identifiers) (pg *types.Page, err error) {
	pg = resource.FindComposePage(rr, ii)
	if pg != nil {
		return pg, nil
	}

	if nsID <= 0 {
		return nil, nil
	}

	// Go in the store
	return findComposePageStore(ctx, s, nsID, makeGenericFilter(ii))
}

// findComposePageStore looks for the page in the store
func findComposePageStore(ctx context.Context, s store.Storer, nsID uint64, gf genericFilter) (pg *types.Page, err error) {
	if nsID == 0 {
		return nil, nil
	}

	if gf.id > 0 {
		pg, err = store.LookupComposePageByID(ctx, s, gf.id)
		if err != nil && err != store.ErrNotFound {
			return nil, err
		}

		if pg != nil {
			return
		}
	}

	for _, i := range gf.identifiers {
		pg, err = store.LookupComposePageByNamespaceIDHandle(ctx, s, nsID, i)
		if err == store.ErrNotFound {
			var pp types.PageSet
			pp, _, err = store.SearchComposePages(ctx, s, types.PageFilter{
				NamespaceID: nsID,
				Title:       i,
				Paging: filter.Paging{
					Limit: 2,
				},
			})
			if len(pp) > 1 {
				return nil, resourceErrIdentifierNotUnique(i)
			}
			if len(pp) == 1 {
				pg = pp[0]
			}
		}

		if err != nil && err != store.ErrNotFound {
			return nil, err
		}

		if pg != nil {
			return
		}
	}

	return nil, nil
}
