package store

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
)

func newAPIGatewayFromResource(res *resource.APIGateway, cfg *EncoderConfig) resourceState {
	return &apiGateway{
		cfg: mergeConfig(cfg, res.Config()),

		res: res,
	}
}

func (n *apiGateway) Prepare(ctx context.Context, pl *payload) (err error) {
	err = n.prepareRoute(ctx, pl)
	if err != nil {
		return err
	}

	return n.prepareFilters(ctx, pl)
}

func (n *apiGateway) prepareRoute(ctx context.Context, pl *payload) (err error) {
	// Reset old identifiers
	n.res.Res.ID = 0

	// Try to get the original workflow
	n.gwr, err = findAPIGatewayStore(ctx, pl.s, makeGenericFilter(n.res.Identifiers()))
	if err != nil {
		return err
	}

	if n.gwr != nil {
		n.res.Res.ID = n.gwr.ID
	}
	return nil
}

func (n *apiGateway) prepareFilters(ctx context.Context, pl *payload) (err error) {
	if n.gwr == nil || n.gwr.ID == 0 {
		return nil
	}

	// Reset old identifiers
	for _, rf := range n.res.Filters {
		rf.Res.ID = 0
	}

	// Try to find any related filters for this route
	tt, _, err := store.SearchApigwFilters(ctx, pl.s, types.ApigwFilterFilter{
		RouteID: n.gwr.ID,
	})
	if err != nil {
		return err
	}

	n.ff = tt
	return nil
}

func (n *apiGateway) Encode(ctx context.Context, pl *payload) (err error) {
	err = n.encodeRoute(ctx, pl)
	if err != nil {
		return err
	}

	return n.encodeFilters(ctx, pl)
}

func (n *apiGateway) encodeRoute(ctx context.Context, pl *payload) (err error) {
	res := n.res.Res
	exists := n.gwr != nil && n.gwr.ID > 0

	// Determine the ID
	if res.ID <= 0 && exists {
		res.ID = n.gwr.ID
	}
	if res.ID <= 0 {
		res.ID = NextID()
	}

	// Sys users
	us, err := resolveUserstamps(ctx, pl.s, pl.state.ParentResources, n.res.Userstamps())
	if err != nil {
		return err
	}

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

	res.CreatedBy = pl.invokerID
	if us != nil {
		if us.CreatedBy != nil {
			res.CreatedBy = us.CreatedBy.UserID
		}
		if us.UpdatedBy != nil {
			res.UpdatedBy = us.UpdatedBy.UserID
		}
		if us.DeletedBy != nil {
			res.DeletedBy = us.DeletedBy.UserID
		}
	}

	// Evaluate the resource skip expression
	// @todo expand available parameters; similar implementation to automation/types/record@Dict
	if skip, err := basicSkipEval(ctx, n.cfg, !exists); err != nil {
		return err
	} else if skip {
		return nil
	}

	// Create a fresh workflow
	if !exists {
		return store.CreateApigwRoute(ctx, pl.s, res)
	}

	// Update existing workflow
	switch n.cfg.OnExisting {
	case resource.Skip:
		return nil

	case resource.MergeLeft:
		res = mergeAPIGateways(n.gwr, res)

	case resource.MergeRight:
		res = mergeAPIGateways(res, n.gwr)
	}

	err = store.UpdateApigwRoute(ctx, pl.s, res)
	if err != nil {
		return err
	}

	n.res.Res = res
	return nil
}

func (n *apiGateway) encodeFilters(ctx context.Context, pl *payload) (err error) {
	exists := len(n.ff) > 0
	ff := make([]*types.ApigwFilter, 0, len(n.res.Filters))

	for _, rf := range n.res.Filters {
		res := rf.Res
		res.Route = n.res.Res.ID
		res.ID = NextID()

		// Sys users
		us, err := resolveUserstamps(ctx, pl.s, pl.state.ParentResources, rf.Userstamps())
		if err != nil {
			return err
		}

		ts := rf.Timestamps()
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
		res.CreatedBy = pl.invokerID
		if us != nil {
			if us.CreatedBy != nil {
				res.CreatedBy = us.CreatedBy.UserID
			}
			if us.UpdatedBy != nil {
				res.UpdatedBy = us.UpdatedBy.UserID
			}
			if us.DeletedBy != nil {
				res.DeletedBy = us.DeletedBy.UserID
			}
		}

		ff = append(ff, res)
	}

	// Create a fresh workflow
	if !exists {
		return store.CreateApigwFilter(ctx, pl.s, ff...)
	}

	// If these filters already exist and we wish to modify them,
	// remove the old ones and create new ones
	switch n.cfg.OnExisting {
	case resource.Skip,
		resource.MergeLeft:
		return nil
	}

	err = store.DeleteApigwFilter(ctx, pl.s, n.ff...)
	return store.CreateApigwFilter(ctx, pl.s, ff...)
}
