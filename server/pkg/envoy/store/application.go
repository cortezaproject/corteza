package store

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	application struct {
		cfg *EncoderConfig

		res *resource.Application
		app *types.Application

		ux *userIndex
	}
)

// mergeApplications merges b into a, prioritising a
func mergeApplications(a, b *types.Application) *types.Application {
	c := *a

	if c.Name == "" {
		c.Name = b.Name
	}
	c.OwnerID = b.OwnerID

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
	if c.Unify == nil || *c.Unify == (types.ApplicationUnify{}) {
		c.Unify = b.Unify
	}

	return &c
}

// findApplication looks for the app in the resources & the store
//
// Provided resources are prioritized.
func findApplication(ctx context.Context, s store.Storer, rr resource.InterfaceSet, ii resource.Identifiers) (ap *types.Application, err error) {
	ap = resource.FindApplication(rr, ii)
	if ap != nil {
		return ap, nil
	}

	return findApplicationStore(ctx, s, makeGenericFilter(ii))
}

// findApplicationStore looks for the app in the store
func findApplicationStore(ctx context.Context, s store.Storer, gf genericFilter) (ap *types.Application, err error) {
	if gf.id > 0 {
		ap, err = store.LookupApplicationByID(ctx, s, gf.id)
		if err != nil && err != store.ErrNotFound {
			return nil, err
		}

		if ap != nil {
			return
		}
	}

	for _, i := range gf.identifiers {
		var aa types.ApplicationSet
		aa, _, err = store.SearchApplications(ctx, s, types.ApplicationFilter{Name: i})
		if err != nil && err != store.ErrNotFound {
			return nil, err
		}
		if len(aa) > 1 {
			return nil, resourceErrIdentifierNotUnique(i)
		}
		if len(aa) == 1 {
			ap = aa[0]
			return
		}
	}

	return nil, nil
}
