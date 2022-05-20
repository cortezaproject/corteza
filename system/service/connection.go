package service

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	a "github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/dal/capabilities"

	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	connection struct {
		actionlog actionlog.Recorder
		store     store.Storer
		ac        connectionAccessController
		dal       dalConnections
	}

	connectionAccessController interface {
		CanGrant(context.Context) bool
		CanSearchConnections(ctx context.Context) bool

		CanCreateConnection(context.Context) bool
		CanReadConnection(context.Context, *types.Connection) bool
		CanUpdateConnection(context.Context, *types.Connection) bool
		CanDeleteConnection(context.Context, *types.Connection) bool
	}

	dalConnections interface {
		AddConnection(ctx context.Context, connectionID uint64, dsn string, dft dal.ConnectionDefaults, capabilities ...capabilities.Capability) (err error)
		UpdateConnection(ctx context.Context, connectionID uint64, dsn string, dft dal.ConnectionDefaults, capabilities ...capabilities.Capability) (err error)
		RemoveConnection(ctx context.Context, connectionID uint64) (err error)
	}
)

func Connection(ctx context.Context, dal dalConnections) (*connection, error) {
	out := &connection{
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		dal:       dal,
	}

	return out, out.reloadConnections(ctx)
}

func (svc *connection) FindByID(ctx context.Context, ID uint64) (q *types.Connection, err error) {
	var (
		rProps = &connectionActionProps{}
	)

	err = func() error {
		if ID == 0 {
			return ConnectionErrInvalidID()
		}

		if q, err = store.LookupConnectionByID(ctx, svc.store, ID); err != nil {
			return ConnectionErrInvalidID().Wrap(err)
		}

		rProps.setConnection(q)

		if !svc.ac.CanReadConnection(ctx, q) {
			return ConnectionErrNotAllowedToRead(rProps)
		}

		return nil
	}()

	return q, svc.recordAction(ctx, rProps, ConnectionActionLookup, err)
}

func (svc *connection) Create(ctx context.Context, new *types.Connection) (q *types.Connection, err error) {
	var (
		qProps = &connectionActionProps{new: new}
	)

	err = func() (err error) {
		if !svc.ac.CanCreateConnection(ctx) {
			return ConnectionErrNotAllowedToCreate(qProps)
		}

		new.ID = nextID()
		new.CreatedAt = *now()
		new.CreatedBy = a.GetIdentityFromContext(ctx).Identity()

		if err = store.CreateConnection(ctx, svc.store, new); err != nil {
			return err
		}

		q = new

		return svc.dal.AddConnection(ctx, new.ID, new.DSN, new.ConnectionDefaults(), new.ActiveCapabilities()...)
	}()

	return q, svc.recordAction(ctx, qProps, ConnectionActionCreate, err)
}

func (svc *connection) Update(ctx context.Context, upd *types.Connection) (q *types.Connection, err error) {
	var (
		qProps = &connectionActionProps{update: upd}
		qq     *types.Connection
		e      error
	)

	err = func() (err error) {
		if qq, e = store.LookupConnectionByID(ctx, svc.store, upd.ID); e != nil {
			return ConnectionErrNotFound(qProps)
		}

		if !svc.ac.CanUpdateConnection(ctx, qq) {
			return ConnectionErrNotAllowedToUpdate(qProps)
		}

		upd.UpdatedAt = now()
		upd.CreatedAt = qq.CreatedAt
		upd.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()

		if err = store.UpdateConnection(ctx, svc.store, upd); err != nil {
			return
		}

		q = upd

		return svc.dal.UpdateConnection(ctx, upd.ID, upd.DSN, upd.ConnectionDefaults(), upd.ActiveCapabilities()...)
	}()

	return q, svc.recordAction(ctx, qProps, ConnectionActionUpdate, err)
}

func (svc *connection) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		qProps = &connectionActionProps{}
		q      *types.Connection
	)

	err = func() (err error) {
		if ID == 0 {
			return ConnectionErrInvalidID()
		}

		if q, err = store.LookupConnectionByID(ctx, svc.store, ID); err != nil {
			return
		}

		if !svc.ac.CanDeleteConnection(ctx, q) {
			return ConnectionErrNotAllowedToDelete(qProps)
		}

		qProps.setConnection(q)

		q.DeletedAt = now()
		q.DeletedBy = a.GetIdentityFromContext(ctx).Identity()

		if err = store.UpdateConnection(ctx, svc.store, q); err != nil {
			return
		}

		return svc.dal.RemoveConnection(ctx, q.ID)
	}()

	return svc.recordAction(ctx, qProps, ConnectionActionDelete, err)
}

func (svc *connection) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		qProps = &connectionActionProps{}
		q      *types.Connection
	)

	err = func() (err error) {
		if ID == 0 {
			return ConnectionErrInvalidID()
		}

		if q, err = store.LookupConnectionByID(ctx, svc.store, ID); err != nil {
			return
		}

		if !svc.ac.CanDeleteConnection(ctx, q) {
			return ConnectionErrNotAllowedToUndelete(qProps)
		}

		qProps.setConnection(q)

		q.DeletedAt = nil
		q.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()

		if err = store.UpdateConnection(ctx, svc.store, q); err != nil {
			return
		}

		return svc.dal.AddConnection(ctx, q.ID, q.DSN, q.ConnectionDefaults(), q.ActiveCapabilities()...)
	}()

	return svc.recordAction(ctx, qProps, ConnectionActionDelete, err)
}

func (svc *connection) Search(ctx context.Context, filter types.ConnectionFilter) (r types.ConnectionSet, f types.ConnectionFilter, err error) {
	var (
		aProps = &connectionActionProps{search: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.Connection) (bool, error) {
		if !svc.ac.CanReadConnection(ctx, res) {
			return false, nil
		}

		return true, nil
	}

	err = func() error {
		if !svc.ac.CanSearchConnections(ctx) {
			return ConnectionErrNotAllowedToSearch()
		}

		if r, f, err = store.SearchConnections(ctx, svc.store, filter); err != nil {
			return err
		}

		return nil
	}()

	return r, f, svc.recordAction(ctx, aProps, ConnectionActionSearch, err)
}

func (svc *connection) reloadConnections(ctx context.Context) (err error) {
	// Get all available connections
	cc, _, err := store.SearchConnections(ctx, svc.store, types.ConnectionFilter{})
	if err != nil {
		return
	}

	for _, c := range cc {
		if err = svc.dal.AddConnection(ctx, c.ID, c.DSN, c.ConnectionDefaults(), c.ActiveCapabilities()...); err != nil {
			return
		}
	}

	return
}
