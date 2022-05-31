package service

import (
	"context"
	"fmt"
	"reflect"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	a "github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/dal/capabilities"
	"github.com/cortezaproject/corteza-server/pkg/options"

	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	dalConnection struct {
		actionlog actionlog.Recorder
		store     store.Storer
		ac        connectionAccessController
		dal       dalConnections
		dbConf    options.DBOpt
	}

	connectionAccessController interface {
		CanGrant(context.Context) bool
		CanSearchDalConnections(ctx context.Context) bool

		CanCreateDalConnection(context.Context) bool
		CanReadDalConnection(context.Context, *types.DalConnection) bool
		CanUpdateDalConnection(context.Context, *types.DalConnection) bool
		CanDeleteDalConnection(context.Context, *types.DalConnection) bool
	}

	dalConnections interface {
		AddConnection(ctx context.Context, connectionID uint64, cp dal.ConnectionParams, dft dal.ConnectionMeta, capabilities ...capabilities.Capability) (err error)
		UpdateConnection(ctx context.Context, connectionID uint64, cp dal.ConnectionParams, dft dal.ConnectionMeta, capabilities ...capabilities.Capability) (err error)
		RemoveConnection(ctx context.Context, connectionID uint64) (err error)
	}
)

func Connection(ctx context.Context, dal dalConnections, dbConf options.DBOpt) *dalConnection {
	return &dalConnection{
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		dal:       dal,
		dbConf:    dbConf,
	}
}

func (svc *dalConnection) FindByID(ctx context.Context, ID uint64) (q *types.DalConnection, err error) {
	var (
		rProps = &dalConnectionActionProps{}
	)

	err = func() error {
		if ID == 0 {
			return DalConnectionErrInvalidID()
		}

		if q, err = store.LookupDalConnectionByID(ctx, svc.store, ID); err != nil {
			return DalConnectionErrInvalidID().Wrap(err)
		}

		rProps.setConnection(q)

		if !svc.ac.CanReadDalConnection(ctx, q) {
			return DalConnectionErrNotAllowedToRead(rProps)
		}

		svc.assurePrimaryConnection(q)
		return nil
	}()
	return q, svc.recordAction(ctx, rProps, DalConnectionActionLookup, err)
}

func (svc *dalConnection) Create(ctx context.Context, new *types.DalConnection) (q *types.DalConnection, err error) {
	var (
		qProps = &dalConnectionActionProps{new: new}
	)

	err = func() (err error) {
		if !svc.ac.CanCreateDalConnection(ctx) {
			return DalConnectionErrNotAllowedToCreate(qProps)
		}

		new.ID = nextID()
		new.CreatedAt = *now()
		new.CreatedBy = a.GetIdentityFromContext(ctx).Identity()

		// validation
		{
			if new.Type != types.DalPrimaryConnectionResourceType {
				// @todo error
				err = fmt.Errorf("cannot create connection: unsupported connection type %s", new.Type)
				return
			}
		}

		if err = store.CreateDalConnection(ctx, svc.store, new); err != nil {
			return err
		}

		q = new

		var cm dal.ConnectionMeta
		cm, err = svc.makeConnectionMeta(ctx, new)
		if err != nil {
			return
		}

		return svc.dal.AddConnection(ctx, new.ID, new.Config.Connection, cm, new.ActiveCapabilities()...)
	}()

	return q, svc.recordAction(ctx, qProps, DalConnectionActionCreate, err)
}

func (svc *dalConnection) Update(ctx context.Context, upd *types.DalConnection) (q *types.DalConnection, err error) {
	var (
		qProps = &dalConnectionActionProps{update: upd}
		old    *types.DalConnection
		e      error
	)

	err = func() (err error) {
		if old, e = store.LookupDalConnectionByID(ctx, svc.store, upd.ID); e != nil {
			return DalConnectionErrNotFound(qProps)
		}

		svc.assurePrimaryConnection(old)

		if !svc.ac.CanUpdateDalConnection(ctx, old) {
			return DalConnectionErrNotAllowedToUpdate(qProps)
		}

		upd.UpdatedAt = now()
		upd.CreatedAt = old.CreatedAt
		upd.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()

		// validate
		{
			if old.Type == types.DalPrimaryConnectionResourceType {
				if !reflect.DeepEqual(old.Config.Connection, upd.Config.Connection) {
					// @todo err
					return fmt.Errorf("can not update connection parameters for primary connection")
				}

				if old.Handle != upd.Handle {
					return fmt.Errorf("can not update handle for primary connection")
				}

				if old.Config.DefaultModelIdent != upd.Config.DefaultModelIdent {
					return fmt.Errorf("can not update defaultModelIdent for primary connection")
				}

				if old.Config.DefaultAttributeIdent != upd.Config.DefaultAttributeIdent {
					return fmt.Errorf("can not update defaultAttributeIdent for primary connection")
				}

				if old.Handle != upd.Handle {
					return fmt.Errorf("can not update handle for primary connection")
				}

				if old.Type != upd.Type {
					return fmt.Errorf("can not update type for primary connection")
				}
			}
		}

		if err = store.UpdateDalConnection(ctx, svc.store, upd); err != nil {
			return
		}

		q = upd
		svc.assurePrimaryConnection(q)

		var cm dal.ConnectionMeta
		cm, err = svc.makeConnectionMeta(ctx, upd)
		if err != nil {
			return
		}

		return svc.dal.UpdateConnection(ctx, upd.ID, upd.Config.Connection, cm, upd.ActiveCapabilities()...)
	}()
	return q, svc.recordAction(ctx, qProps, DalConnectionActionUpdate, err)
}

func (svc *dalConnection) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		qProps = &dalConnectionActionProps{}
		q      *types.DalConnection
	)

	err = func() (err error) {
		if ID == 0 {
			return DalConnectionErrInvalidID()
		}

		if q, err = store.LookupDalConnectionByID(ctx, svc.store, ID); err != nil {
			return
		}

		if q.Type == types.DalPrimaryConnectionResourceType {
			return fmt.Errorf("not allowed to delete primary connections")
		}

		if !svc.ac.CanDeleteDalConnection(ctx, q) {
			return DalConnectionErrNotAllowedToDelete(qProps)
		}

		qProps.setConnection(q)

		q.DeletedAt = now()
		q.DeletedBy = a.GetIdentityFromContext(ctx).Identity()

		if err = store.UpdateDalConnection(ctx, svc.store, q); err != nil {
			return
		}

		return svc.dal.RemoveConnection(ctx, q.ID)
	}()

	return svc.recordAction(ctx, qProps, DalConnectionActionDelete, err)
}

func (svc *dalConnection) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		qProps = &dalConnectionActionProps{}
		q      *types.DalConnection
	)

	err = func() (err error) {
		if ID == 0 {
			return DalConnectionErrInvalidID()
		}

		if q, err = store.LookupDalConnectionByID(ctx, svc.store, ID); err != nil {
			return
		}

		if !svc.ac.CanDeleteDalConnection(ctx, q) {
			return DalConnectionErrNotAllowedToUndelete(qProps)
		}

		qProps.setConnection(q)

		q.DeletedAt = nil
		q.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()

		if err = store.UpdateDalConnection(ctx, svc.store, q); err != nil {
			return
		}

		var cm dal.ConnectionMeta
		cm, err = svc.makeConnectionMeta(ctx, q)
		if err != nil {
			return
		}

		return svc.dal.AddConnection(ctx, q.ID, q.Config.Connection, cm, q.ActiveCapabilities()...)
	}()

	return svc.recordAction(ctx, qProps, DalConnectionActionDelete, err)
}

func (svc *dalConnection) Search(ctx context.Context, filter types.DalConnectionFilter) (r types.DalConnectionSet, f types.DalConnectionFilter, err error) {
	var (
		aProps = &dalConnectionActionProps{search: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.DalConnection) (bool, error) {
		if !svc.ac.CanReadDalConnection(ctx, res) {
			return false, nil
		}

		return true, nil
	}

	err = func() error {
		if !svc.ac.CanSearchDalConnections(ctx) {
			return DalConnectionErrNotAllowedToSearch()
		}

		if r, f, err = store.SearchDalConnections(ctx, svc.store, filter); err != nil {
			return err
		}

		svc.assurePrimaryConnection(r...)
		return nil
	}()
	return r, f, svc.recordAction(ctx, aProps, DalConnectionActionSearch, err)
}

func (svc *dalConnection) ReloadConnections(ctx context.Context) (err error) {
	// Get all available connections
	cc, _, err := store.SearchDalConnections(ctx, svc.store, types.DalConnectionFilter{})
	if err != nil {
		return
	}

	for _, c := range cc {
		var cm dal.ConnectionMeta
		cm, err = svc.makeConnectionMeta(ctx, c)
		if err != nil {
			return
		}
		if err = svc.dal.AddConnection(ctx, c.ID, c.Config.Connection, cm, c.ActiveCapabilities()...); err != nil {
			return
		}
	}

	return
}

func (svc *dalConnection) makeConnectionMeta(ctx context.Context, c *types.DalConnection) (cm dal.ConnectionMeta, err error) {
	// @todo we could probably utilize connection params more here
	cm = dal.ConnectionMeta{
		DefaultModelIdent:      c.Config.DefaultModelIdent,
		DefaultAttributeIdent:  c.Config.DefaultAttributeIdent,
		DefaultPartitionFormat: c.Config.DefaultPartitionFormat,
		SensitivityLevel:       c.SensitivityLevel,
		Label:                  c.Handle,
	}

	return
}

func (svc *dalConnection) assurePrimaryConnection(connections ...*types.DalConnection) {
	for _, c := range connections {
		if c.Type == types.DalPrimaryConnectionResourceType {
			c.Config.Connection = dal.NewDSNConnection(svc.dbConf.DSN)
			return
		}
	}
}
