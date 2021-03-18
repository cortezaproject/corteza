package store

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/store"
)

type (
	composeModule struct {
		cfg *EncoderConfig

		res *resource.ComposeModule
		mod *types.Module

		relNS     *types.Namespace
		recFields map[string]uint64
	}
)

// mergeComposeModuleFields merges b into a, prioritising a
func mergeComposeModuleFields(a, b types.ModuleFieldSet) types.ModuleFieldSet {
	ff := a.Clone()
	missing := make(types.ModuleFieldSet, 0, len(b))

	for _, fb := range b {
		for _, fa := range ff {
			if fb.ID != fa.ID && fb.Name != fa.Name {
				continue
			}

			fa.ModuleID = fb.ModuleID
			fa.Place = fb.Place
			if fa.Kind == "" {
				fa.Kind = fb.Kind
			}
			if fa.Name == "" {
				fa.Name = fb.Name
			}
			if fa.Label == "" {
				fa.Label = fb.Label
			}
			if fa.Options == nil {
				fa.Options = fb.Options
			}
			if fa.DefaultValue == nil {
				fa.DefaultValue = fb.DefaultValue
			}
			if fa.Expressions.ValueExpr == "" {
				fa.Expressions.ValueExpr = fb.Expressions.ValueExpr
			}
			if len(fa.Expressions.Sanitizers) == 0 {
				fa.Expressions.Sanitizers = fb.Expressions.Sanitizers
			}
			if len(fa.Expressions.Validators) == 0 {
				fa.Expressions.Validators = fb.Expressions.Validators
			}
			if len(fa.Expressions.Formatters) == 0 {
				fa.Expressions.Formatters = fb.Expressions.Formatters
			}

			if fa.CreatedAt.IsZero() {
				fa.CreatedAt = fb.CreatedAt
			}
			if fa.UpdatedAt == nil {
				fa.UpdatedAt = fb.UpdatedAt
			}
			if fa.DeletedAt == nil {
				fa.DeletedAt = fb.DeletedAt
			}

			goto out
		}
		missing = append(missing, fb)
	out:
	}

	ff = append(ff, missing...)

	return ff
}

// mergeComposeModule merges b into a, prioritising a
func mergeComposeModule(a, b *types.Module) *types.Module {
	c := a.Clone()

	if c.Handle == "" {
		c.Handle = b.Handle
	}
	if c.Name == "" {
		c.Name = b.Name
	}
	c.NamespaceID = b.NamespaceID

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
	if c.Meta == nil {
		c.Meta = b.Meta
	}

	return c
}

// findComposeModule looks for the chart in the resources & the store
//
// Provided resources are prioritized.
func findComposeModule(ctx context.Context, s store.Storer, nsID uint64, rr resource.InterfaceSet, ii resource.Identifiers) (mod *types.Module, err error) {
	mod = resource.FindComposeModule(rr, ii)
	if mod != nil {
		return mod, nil
	}

	if nsID <= 0 {
		return nil, nil
	}

	// Go in the store
	return findComposeModuleStore(ctx, s, nsID, makeGenericFilter(ii))
}

// findComposeModuleStore looks for the module in the store
func findComposeModuleStore(ctx context.Context, s store.Storer, nsID uint64, gf genericFilter) (mod *types.Module, err error) {
	if nsID == 0 {
		return nil, nil
	}

	if gf.id > 0 {
		mod, err = store.LookupComposeModuleByID(ctx, s, gf.id)
		if err != nil && err != store.ErrNotFound {
			return nil, err
		}

		if mod != nil {
			return
		}
	}

	for _, i := range gf.identifiers {
		mod, err = store.LookupComposeModuleByNamespaceIDHandle(ctx, s, nsID, i)
		if err == store.ErrNotFound {
			var mm types.ModuleSet
			mm, _, err = store.SearchComposeModules(ctx, s, types.ModuleFilter{
				NamespaceID: nsID,
				Name:        i,
				Paging: filter.Paging{
					Limit: 2,
				},
			})
			if len(mm) > 1 {
				return nil, resourceErrIdentifierNotUnique(i)
			}
			if len(mm) == 1 {
				mod = mm[0]
			}
		}

		if err != nil && err != store.ErrNotFound {
			return nil, err
		}

		if mod != nil {
			return
		}
	}

	return nil, nil
}

// findComposeModuleFieldsStore looks for the module fields in the store
func findComposeModuleFieldsStore(ctx context.Context, s store.Storer, mod *types.Module) (types.ModuleFieldSet, error) {
	if mod.ID <= 0 {
		return mod.Fields, nil
	}

	ff, _, err := store.SearchComposeModuleFields(ctx, s, types.ModuleFieldFilter{
		ModuleID: []uint64{mod.ID},
	})
	if err != nil {
		return nil, err
	}

	return ff, nil
}

func composeModuleErrUnresolvedRecordField(ii resource.Identifiers) error {
	return fmt.Errorf("record module field unresolved %v", ii.StringSlice())
}
