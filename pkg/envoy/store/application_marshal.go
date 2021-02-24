package store

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
)

func NewApplicationFromResource(res *resource.Application, cfg *EncoderConfig) resourceState {
	return &application{
		cfg: mergeConfig(cfg, res.Config()),
		res: res,
	}
}

// Prepare prepares the application to be encoded
//
// Any validation, additional constraining should be performed here.
func (n *application) Prepare(ctx context.Context, pl *payload) (err error) {
	// Get the existing app
	n.app, err = findApplicationS(ctx, pl.s, makeGenericFilter(n.res.Identifiers()))
	if err != nil {
		return err
	}

	if n.app != nil {
		n.res.Res.ID = n.app.ID
	}
	return nil
}

// Encode encodes the application to the store
//
// Encode is allowed to do some data manipulation, but no resource constraints
// should be changed.
func (n *application) Encode(ctx context.Context, pl *payload) (err error) {
	res := n.res.Res
	exists := n.app != nil && n.app.ID > 0

	// Determine the ID
	if res.ID <= 0 && exists {
		res.ID = n.app.ID
	}
	if res.ID <= 0 {
		res.ID = NextID()
	}

	// Sys users
	us, err := resolveUserstamps(ctx, pl.s, pl.state.ParentResources, n.res.Userstamps())
	if err != nil {
		return err
	}

	// Unify
	if res.Unify == nil {
		res.Unify = &types.ApplicationUnify{}
	}

	// Timestamps
	ts := n.res.Timestamps()
	if ts != nil {
		if ts.CreatedAt != nil {
			res.CreatedAt = *ts.CreatedAt.T
		} else {
			res.CreatedAt = *now()
		}
		if ts.UpdatedAt != nil {
			res.UpdatedAt = ts.UpdatedAt.T
		}
		if ts.DeletedAt != nil {
			res.DeletedAt = ts.DeletedAt.T
		}
	}

	// Userstamps
	if us != nil {
		if us.OwnedBy != nil {
			res.OwnerID = us.OwnedBy.UserID
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
		return store.CreateApplication(ctx, pl.s, res)
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

	err = store.UpdateApplication(ctx, pl.s, res)
	if err != nil {
		return err
	}

	n.res.Res = res
	return nil
}
