package store

import (
	"context"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/store"
)

type (
	composeNamespace struct {
		cfg *EncoderConfig

		res *resource.ComposeNamespace
		ns  *types.Namespace
	}

	composeNamespaceSet []*composeNamespace
)

// mergeComposeNamespaces merges b into a, prioritising a
func mergeComposeNamespaces(a, b *types.Namespace) *types.Namespace {
	c := a.Clone()

	if c.Name == "" {
		c.Name = b.Name
	}
	if c.Slug == "" {
		c.Slug = b.Slug
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

	// I'll just compare the entire struct for now
	if c.Meta == (types.NamespaceMeta{}) {
		c.Meta = b.Meta
	}

	return c
}

// findComposeNamespace looks for the namespace in the resources & the store
//
// Provided resources are prioritized.
func findComposeNamespace(ctx context.Context, s store.Storer, rr resource.InterfaceSet, ii resource.Identifiers) (ns *types.Namespace, err error) {
	ns = resource.FindComposeNamespace(rr, ii)
	if ns != nil {
		return ns, nil
	}

	return findComposeNamespaceStore(ctx, s, makeGenericFilter(ii))
}

// findComposeNamespaceStore looks for the namespace in the store
func findComposeNamespaceStore(ctx context.Context, s store.Storer, gf genericFilter) (ns *types.Namespace, err error) {
	if gf.id > 0 {
		ns, err = store.LookupComposeNamespaceByID(ctx, s, gf.id)
		if err != nil && err != store.ErrNotFound {
			return nil, err
		}

		if ns != nil {
			return
		}
	}

	for _, i := range gf.identifiers {
		ns, err = store.LookupComposeNamespaceBySlug(ctx, s, i)
		if err == store.ErrNotFound {
			var nn types.NamespaceSet
			nn, _, err = store.SearchComposeNamespaces(ctx, s, types.NamespaceFilter{
				Name: i,
				Paging: filter.Paging{
					Limit: 2,
				},
			})
			if len(nn) > 1 {
				return nil, resourceErrIdentifierNotUnique(i)
			}
			if len(nn) == 1 {
				ns = nn[0]
			}
		}

		if err != nil && err != store.ErrNotFound {
			return nil, err
		}

		if ns != nil {
			return
		}
	}

	return nil, nil
}
