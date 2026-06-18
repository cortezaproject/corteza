package service

import (
	"context"

	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/pkg/apigw"
	a "github.com/crusttech/human/server/pkg/auth"

	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

// The CRUD skeleton (FindByID, Search, Create, Update, DeleteByID,
// UndeleteByID and loadApigwRoute) is generated in apigw_route.gen.go from
// system/apigw_route.cue.
//
// This file owns the struct, access-controller interface, constructor and the
// custom on<Op> bodies that the generated Create / Update / DeleteByID /
// UndeleteByID delegate to (CreatedBy/Group defaulting, endpoint-moved 404
// handling, apigw reload / not-found signalling, soft-delete via
// UpdateApigwRoute).

type (
	apigwRoute struct {
		actionlog actionlog.Recorder
		store     store.Storer
		ac        routeAccessController
	}

	routeAccessController interface {
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

// onCreate is the custom body for the generated Create. The generated method
// owns the action-log scaffold + recordAction + the standard CanCreateApigwRoute
// check; everything below (CreatedBy/Group defaulting, store create, apigw
// reload signalling) lives here.
func (svc *apigwRoute) onCreate(ctx context.Context, new *types.ApigwRoute) (err error) {
	new.ID = nextID()
	new.CreatedAt = *now()
	new.CreatedBy = a.GetIdentityFromContext(ctx).Identity()

	// todo
	new.Group = 0

	if err = store.CreateApigwRoute(ctx, svc.store, new); err != nil {
		return err
	}

	// send the signal to reload new route
	if new.Enabled {
		if err = apigw.Service().ReloadEndpoint(ctx, new.Method, new.Endpoint); err != nil {
			return err
		}
	}

	return nil
}

// onUpdate is the custom body for the generated Update. The generated method
// owns the action-log scaffold + recordAction; everything below (load,
// endpoint-moved detection, resource-based access, stale check, field copy,
// store update, apigw not-found / reload signalling) lives here.
func (svc *apigwRoute) onUpdate(ctx context.Context, upd *types.ApigwRoute, qProps *apigwRouteActionProps) (res *types.ApigwRoute, err error) {
	var e error

	if res, e = loadApigwRoute(ctx, svc.store, upd.ID); e != nil {
		return res, ApigwRouteErrNotFound(qProps)
	}

	var (
		// check if old endpoint moved and attach the 404 handler
		endpointMoved = res.Enabled != upd.Enabled || res.Method != upd.Method || res.Endpoint != upd.Endpoint
	)

	if !svc.ac.CanUpdateApigwRoute(ctx, res) {
		return res, ApigwRouteErrNotAllowedToUpdate(qProps)
	}

	// Test if stale (update has an older version of data)
	if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
		return res, ApigwRouteErrStaleData()
	}

	// copy (potentially) updated files from the payload
	res.Meta = upd.Meta
	res.Method = upd.Method
	res.Endpoint = upd.Endpoint
	res.Enabled = upd.Enabled
	res.Group = upd.Group

	// ensure we have a valid endpoint
	res.UpdatedAt = now()
	res.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()

	if err = store.UpdateApigwRoute(ctx, svc.store, res); err != nil {
		return res, err
	}

	// @todo move this into struct of the service (svc),
	//       same as we do for services for other structs
	ags := apigw.Service()

	// old endpoint moved: attach 404 handler
	if endpointMoved {
		ags.NotFound(ctx, res.Method, res.Endpoint)
	}

	// send the signal to reload updated route
	if upd.Enabled {
		if err = ags.ReloadEndpoint(ctx, upd.Method, upd.Endpoint); err != nil {
			return res, err
		}
	}

	return res, nil
}

// onDelete is the custom body for the generated DeleteByID. The generated method
// owns the action-log scaffold + recordAction; everything below (load,
// resource-based access, soft-delete via UpdateApigwRoute, apigw not-found
// signalling) lives here.
func (svc *apigwRoute) onDelete(ctx context.Context, ID uint64, qProps *apigwRouteActionProps) (err error) {
	var q *types.ApigwRoute

	if q, err = loadApigwRoute(ctx, svc.store, ID); err != nil {
		return
	}

	if !svc.ac.CanDeleteApigwRoute(ctx, q) {
		return ApigwRouteErrNotAllowedToDelete(qProps)
	}

	qProps.setRoute(q)

	q.DeletedAt = now()
	q.DeletedBy = a.GetIdentityFromContext(ctx).Identity()

	if err = store.UpdateApigwRoute(ctx, svc.store, q); err != nil {
		return
	}

	// send the signal to reload deleted route
	if q.Enabled {
		apigw.Service().NotFound(ctx, q.Method, q.Endpoint)
	}

	return nil
}

// onUndelete is the custom body for the generated UndeleteByID. The generated
// method owns the action-log scaffold + recordAction; everything below (load,
// resource-based access, clear deleted_at, apigw reload signalling) lives here.
func (svc *apigwRoute) onUndelete(ctx context.Context, ID uint64, qProps *apigwRouteActionProps) (err error) {
	var q *types.ApigwRoute

	if q, err = loadApigwRoute(ctx, svc.store, ID); err != nil {
		return
	}

	if !svc.ac.CanDeleteApigwRoute(ctx, q) {
		return ApigwRouteErrNotAllowedToUndelete(qProps)
	}

	qProps.setRoute(q)

	q.DeletedAt = nil
	q.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()

	if err = store.UpdateApigwRoute(ctx, svc.store, q); err != nil {
		return
	}

	// send the signal to reload all queues
	if q.Enabled {
		if err = apigw.Service().ReloadEndpoint(ctx, q.Method, q.Endpoint); err != nil {
			return err
		}
	}

	return nil
}
