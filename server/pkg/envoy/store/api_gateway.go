package store

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	apiGateway struct {
		cfg *EncoderConfig

		res *resource.APIGateway
		gwr *types.ApigwRoute
		ff  types.ApigwFilterSet

		ux *userIndex
	}
	apiGatewaySet []*apiGateway

	apiGatewayFilter struct {
		cfg *EncoderConfig

		res *resource.APIGatewayFilter
		tr  *types.ApigwFilter
	}
	apiGatewayFilterSet []*apiGatewayFilter
)

// mergeAPIGateways merges b into a, prioritising a
func mergeAPIGateways(a, b *types.ApigwRoute) *types.ApigwRoute {
	c := a

	if c.Endpoint == "" {
		c.Endpoint = b.Endpoint
	}
	if c.Method == "" {
		c.Method = b.Method
	}

	c.Enabled = b.Enabled

	if c.Group == 0 {
		c.Group = b.Group
	}

	c.Meta = b.Meta

	if c.CreatedBy == 0 {
		c.CreatedBy = b.CreatedBy
	}
	if c.UpdatedBy == 0 {
		c.UpdatedBy = b.UpdatedBy
	}
	if c.DeletedBy == 0 {
		c.DeletedBy = b.DeletedBy
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

// findAPIGateway looks for the workflow in the resources & the store
//
// Provided resources are prioritized.
func findAPIGateway(ctx context.Context, s store.Storer, rr resource.InterfaceSet, ii resource.Identifiers) (wf *types.ApigwRoute, err error) {
	wf = resource.FindAPIGateway(rr, ii)
	if wf != nil {
		return wf, nil
	}

	return findAPIGatewayStore(ctx, s, makeGenericFilter(ii))
}

// findAPIGatewayStore looks for the workflow in the store
func findAPIGatewayStore(ctx context.Context, s store.Storer, gf genericFilter) (wf *types.ApigwRoute, err error) {
	if gf.id > 0 {
		wf, err = store.LookupApigwRouteByID(ctx, s, gf.id)
		if err != nil && err != store.ErrNotFound {
			return nil, err
		}

		if wf != nil {
			return
		}
	}

	for _, i := range gf.identifiers {
		wf, err = store.LookupApigwRouteByEndpoint(ctx, s, i)
		if err != nil && err != store.ErrNotFound {
			return nil, err
		}

		if wf != nil {
			return
		}
	}

	return nil, nil
}
