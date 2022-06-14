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
	"github.com/cortezaproject/corteza-server/system/dalutils"
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
		CreateConnection(ctx context.Context, connectionID uint64, cp dal.ConnectionParams, dft dal.ConnectionMeta, capabilities ...capabilities.Capability) (err error)
		UpdateConnection(ctx context.Context, connectionID uint64, cp dal.ConnectionParams, dft dal.ConnectionMeta, capabilities ...capabilities.Capability) (err error)
		DeleteConnection(ctx context.Context, connectionID uint64) (err error)
		SearchConnectionIssues(connectionID uint64) (err []error)
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

		svc.proc(q)
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
			if new.Type != types.DalConnectionResourceType {
				// @todo error
				err = fmt.Errorf("cannot create connection: unsupported connection type %s", new.Type)
				return
			}
		}

		if err = store.CreateDalConnection(ctx, svc.store, new); err != nil {
			return err
		}

		q = new

		if err = dalutils.DalConnectionCreate(ctx, svc.dal, new); err != nil {
			return err
		}
		svc.proc(q)
		return
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

		svc.proc(old)

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
		defer svc.proc(q)
		if old.HasIssues() {
			return dalutils.DalConnectionCreate(ctx, svc.dal, upd)
		}

		return dalutils.DalConnectionUpdate(ctx, svc.dal, upd)
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

		return dalutils.DalConnectionDelete(ctx, svc.dal, q)
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

		// We're creating it here since it was removed on delete
		return dalutils.DalConnectionCreate(ctx, svc.dal, q)
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

		svc.proc(r...)
		return nil
	}()
	return r, f, svc.recordAction(ctx, aProps, DalConnectionActionSearch, err)
}

func (svc *dalConnection) ReloadConnections(ctx context.Context) (err error) {
	return dalutils.DalConnectionReload(ctx, svc.store, svc.dal)
}

func (svc *dalConnection) proc(connections ...*types.DalConnection) {
	for _, c := range connections {
		svc.procPrimaryConnection(c)
		svc.procDal(c)
		svc.procLocale(c)
	}
}

func (svc *dalConnection) procPrimaryConnection(c *types.DalConnection) {
	if c.Type == types.DalPrimaryConnectionResourceType {
		c.Config.Connection = dal.NewDSNConnection(svc.dbConf.DSN)
		return
	}
}

func (svc *dalConnection) procDal(c *types.DalConnection) {
	ii := svc.dal.SearchConnectionIssues(c.ID)
	if len(ii) == 0 {
		c.Issues = nil
		return
	}

	c.Issues = make([]string, len(ii))
	for i, err := range ii {
		c.Issues[i] = err.Error()
	}
}

func (svc *dalConnection) procLocale(c *types.DalConnection) {
	// @todo...
}
