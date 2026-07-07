package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"

	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/pkg/scope"
	"github.com/crusttech/human/server/store"
	types "github.com/crusttech/human/server/system/types"
)

type connectionServices struct {
	scope scope.Scope
	caps  scope.Capabilities
}

func (svc *connection) Create(ctx context.Context, new *types.Connection) (res *types.Connection, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &connectionActionProps{connection: new, new: new}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}
		if !svc.ac.CanCreateConnection(ctx) {
			return ConnectionErrNotAllowedToCreate()
		}

		if err = svc.beforeCreate(ctx, new); err != nil {
			return err
		}

		new.ID = nextID()
		new.CreatedAt = *now()

		if err = store.CreateConnection(ctx, svc.store, new); err != nil {
			return
		}

		if err = label.Create(ctx, svc.store, new); err != nil {
			return
		}

		res = new

		if err = svc.afterCreate(ctx, res); err != nil {
			return err
		}
		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ConnectionActionCreate, err)
}

func (svc *connection) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &connectionActionProps{}
		res    *types.Connection
	)
	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		if res, err = loadConnection(ctx, svc.store, ID); err != nil {
			return
		}

		aProps.setConnection(res)

		if !svc.ac.CanDeleteConnection(ctx, res) {
			return ConnectionErrNotAllowedToDelete()
		}

		if err = svc.beforeDelete(ctx, res); err != nil {
			return err
		}

		res.DeletedAt = now()
		if err = store.UpdateConnection(ctx, svc.store, res); err != nil {
			return
		}
		return nil
	}()

	return svc.recordAction(ctx, aProps, ConnectionActionDelete, err)
}

func (svc *connection) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &connectionActionProps{}
		res    *types.Connection
	)
	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		if res, err = loadConnection(ctx, svc.store, ID); err != nil {
			return
		}

		aProps.setConnection(res)

		if !svc.ac.CanDeleteConnection(ctx, res) {
			return ConnectionErrNotAllowedToUndelete()
		}

		if err = svc.beforeUndelete(ctx, res); err != nil {
			return err
		}

		res.DeletedAt = nil
		if err = store.UpdateConnection(ctx, svc.store, res); err != nil {
			return
		}

		return nil
	}()

	return svc.recordAction(ctx, aProps, ConnectionActionUndelete, err)
}

func loadConnection(ctx context.Context, s store.Connections, ID uint64) (res *types.Connection, err error) {
	if ID == 0 {
		return nil, ConnectionErrInvalidID()
	}

	if res, err = store.LookupConnectionByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, ConnectionErrNotFound()
	}

	return
}

// toLabeledConnections converts to []label.LabeledResource
func toLabeledConnections(set []*types.Connection) []label.LabeledResource {
	if len(set) == 0 {
		return nil
	}

	ll := make([]label.LabeledResource, len(set))
	for i := range set {
		ll[i] = set[i]
	}

	return ll
}

func (svc *connection) checkScope(ctx context.Context, cap scope.Capability) error {
	if err := scope.RequireTenantMembership(ctx); err != nil {
		return err
	}
	return scope.RequireCapability(ctx, cap)
}

func (svc *connection) scopeServices(ctx context.Context) *connectionServices {
	return &connectionServices{
		scope: scope.GetScopeFromContext(ctx),
		caps:  scope.GetCapabilitiesFromContext(ctx),
	}
}
