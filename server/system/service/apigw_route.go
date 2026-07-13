package service

import (
	"context"

	"github.com/crusttech/human/server/pkg/apigw"
	a "github.com/crusttech/human/server/pkg/auth"

	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

type (
	apigwRouteAccessController interface {
		CanGrant(context.Context) bool
		CanSearchApigwRoutes(ctx context.Context) bool

		CanCreateApigwRoute(context.Context) bool
		CanReadApigwRoute(context.Context, *types.ApigwRoute) bool
		CanUpdateApigwRoute(context.Context, *types.ApigwRoute) bool
		CanDeleteApigwRoute(context.Context, *types.ApigwRoute) bool
	}
)

func Route() *apigwRoute {
	return &apigwRoute{
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
	}
}

func (svc *apigwRoute) onCreate(ctx context.Context, new *types.ApigwRoute) error {
	new.ID = nextID()
	new.CreatedAt = *now()
	new.CreatedBy = a.GetIdentityFromContext(ctx).Identity()
	new.Group = 0

	if err := store.CreateApigwRoute(ctx, svc.store, new); err != nil {
		return err
	}

	if new.Enabled {
		return apigw.Service().ReloadEndpoint(ctx, new.Method, new.Endpoint)
	}

	return nil
}

func (svc *apigwRoute) onUpdate(ctx context.Context, s store.Storer, upd *types.ApigwRoute, res *types.ApigwRoute, aProps *apigwRouteActionProps, before func() error, after func() error) error {
	if !svc.ac.CanUpdateApigwRoute(ctx, res) {
		return ApigwRouteErrNotAllowedToUpdate(aProps)
	}

	// res still holds old values here; gen copies Endpoint/Method/Enabled/Group after this hook.
	endpointMoved := res.Enabled != upd.Enabled || res.Method != upd.Method || res.Endpoint != upd.Endpoint

	// Meta is not copied by gen.
	res.Meta = upd.Meta
	upd.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()

	ags := apigw.Service()

	if endpointMoved {
		ags.NotFound(ctx, res.Method, res.Endpoint)
	}

	if upd.Enabled {
		return ags.ReloadEndpoint(ctx, upd.Method, upd.Endpoint)
	}

	return nil
}

func (svc *apigwRoute) onDelete(ctx context.Context, s store.Storer, res *types.ApigwRoute, aProps *apigwRouteActionProps) error {
	if !svc.ac.CanDeleteApigwRoute(ctx, res) {
		return ApigwRouteErrNotAllowedToDelete(aProps)
	}

	res.DeletedAt = now()
	res.DeletedBy = a.GetIdentityFromContext(ctx).Identity()

	if err := store.UpdateApigwRoute(ctx, s, res); err != nil {
		return err
	}

	if res.Enabled {
		apigw.Service().NotFound(ctx, res.Method, res.Endpoint)
	}

	return nil
}

func (svc *apigwRoute) onUndelete(ctx context.Context, s store.Storer, res *types.ApigwRoute, aProps *apigwRouteActionProps) error {
	if !svc.ac.CanDeleteApigwRoute(ctx, res) {
		return ApigwRouteErrNotAllowedToUndelete(aProps)
	}

	res.DeletedAt = nil
	res.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()

	if err := store.UpdateApigwRoute(ctx, s, res); err != nil {
		return err
	}

	if res.Enabled {
		return apigw.Service().ReloadEndpoint(ctx, res.Method, res.Endpoint)
	}

	return nil
}
