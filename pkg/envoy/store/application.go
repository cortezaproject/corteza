package store

import (
	"context"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	applicationState struct {
		cfg *EncoderConfig

		res *resource.Application
		app *types.Application
	}
)

func NewApplicationState(res *resource.Application, cfg *EncoderConfig) resourceState {
	return &applicationState{
		cfg: mergeConfig(cfg, res.Config()),
		res: res,
	}
}

func (n *applicationState) Prepare(ctx context.Context, s store.Storer, rs *envoy.ResourceState) (err error) {
	// Initial values
	if n.res.Res.CreatedAt.IsZero() {
		n.res.Res.CreatedAt = time.Now()
	}
	if n.res.Res.Unify == nil {
		n.res.Res.Unify = &types.ApplicationUnify{}
	}

	// Get the existing app
	n.app, err = findApplicationS(ctx, s, makeGenericFilter(n.res.Identifiers()))
	if err != nil {
		return err
	}

	if n.app != nil {
		n.res.Res.ID = n.app.ID
	}
	return nil
}

// Encode encodes the given application
func (n *applicationState) Encode(ctx context.Context, s store.Storer, state *envoy.ResourceState) (err error) {
	res := n.res.Res
	exists := n.app != nil && n.app.ID > 0

	// Determine the ID
	if res.ID <= 0 && exists {
		res.ID = n.app.ID
	}
	if res.ID <= 0 {
		res.ID = NextID()
	}

	// This is not possible, but let's do it anyway
	if state.Conflicting {
		return nil
	}

	// Timestamps
	ts := n.res.Timestamps()
	if ts != nil {
		if ts.CreatedAt != "" {
			t := toTime(ts.CreatedAt)
			if t != nil {
				res.CreatedAt = *t
			}
		}
		if ts.UpdatedAt != "" {
			res.UpdatedAt = toTime(ts.UpdatedAt)
		}
		if ts.DeletedAt != "" {
			res.DeletedAt = toTime(ts.DeletedAt)
		}
	}

	// Evaluate the resource skip expression
	// @todo expand available parameters; similar implementation to compose/types/record@Dict
	if skip, err := basicSkipEval(ctx, n.cfg, !exists); err != nil {
		return err
	} else if skip {
		return nil
	}

	// Create fresh application
	if !exists {
		return store.CreateApplication(ctx, s, res)
	}

	// Update existing application
	switch n.cfg.OnExisting {
	case resource.Skip:
		return nil

	case resource.MergeLeft:
		res = mergeApplications(n.app, res)

	case resource.MergeRight:
		res = mergeApplications(res, n.app)
	}

	err = store.UpdateApplication(ctx, s, res)
	if err != nil {
		return err
	}

	n.res.Res = res
	return nil
}

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

// findApplicationRS looks for the app in the resources & the store
//
// Provided resources are prioritized.
func findApplicationRS(ctx context.Context, s store.Storer, rr resource.InterfaceSet, ii resource.Identifiers) (ap *types.Application, err error) {
	ap = findApplicationR(rr, ii)
	if ap != nil {
		return ap, nil
	}

	return findApplicationS(ctx, s, makeGenericFilter(ii))
}

// findApplicationS looks for the app in the store
func findApplicationS(ctx context.Context, s store.Storer, gf genericFilter) (ap *types.Application, err error) {
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

// findApplicationR looks for the app in the resource set
func findApplicationR(rr resource.InterfaceSet, ii resource.Identifiers) (ap *types.Application) {
	var apRes *resource.Application

	rr.Walk(func(r resource.Interface) error {
		ar, ok := r.(*resource.Application)
		if !ok {
			return nil
		}

		if ar.Identifiers().HasAny(ii) {
			apRes = ar
		}

		return nil
	})

	// Found it
	if apRes != nil {
		return apRes.Res
	}

	return nil
}
