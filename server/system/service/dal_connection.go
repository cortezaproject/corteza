package service

import (
	"context"
	"fmt"
	"regexp"

	"github.com/crusttech/human/server/pkg/errors"

	"github.com/crusttech/human/server/pkg/actionlog"
	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/dal"
	"github.com/crusttech/human/server/pkg/options"

	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

type (
	dalConnection struct {
		actionlog actionlog.Recorder
		store     store.Storer
		ac        dalConnectionAccessController
		dal       dalConnManager
		dbConf    options.DBOpt
	}

	dalConnectionAccessController interface {
		CanGrant(context.Context) bool
		CanSearchDalConnections(ctx context.Context) bool

		CanCreateDalConnection(context.Context) bool
		CanReadDalConnection(context.Context, *types.DalConnection) bool
		CanUpdateDalConnection(context.Context, *types.DalConnection) bool
		CanDeleteDalConnection(context.Context, *types.DalConnection) bool
		CanManageDalConfigOnDalConnection(context.Context, *types.DalConnection) bool
	}

	// Connection management on DAL Service
	dalConnManager interface {
		ReplaceConnection(context.Context, *dal.ConnectionWrap, bool) error
		RemoveConnection(context.Context, uint64) error
		SearchConnectionIssues(uint64) []dal.Issue
	}
)

func DalConnection(ctx context.Context, dal dalConnManager, dbConf options.DBOpt) *dalConnection {
	return &dalConnection{
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		dal:       dal,
		dbConf:    dbConf,
	}
}

// onLookup is the custom body for the generated FindByID. The generated method
// owns the action-log scaffold + recordAction; the load-error wrapping, access
// check (CanReadDalConnection) and proc enrichment, which the standard template
// does not emit, live here.
func (svc *dalConnection) onLookup(ctx context.Context, ID uint64, aProps *dalConnectionActionProps) (q *types.DalConnection, err error) {
	if q, err = loadDalConnection(ctx, svc.store, ID); err != nil {
		return nil, DalConnectionErrInvalidID().Wrap(err)
	}

	aProps.setConnection(q)

	if !svc.ac.CanReadDalConnection(ctx, q) {
		return nil, DalConnectionErrNotAllowedToRead(aProps)
	}

	svc.proc(ctx, q)
	return q, nil
}

// onSearch is the custom body for the generated Search. The generated method
// owns the action-log scaffold + recordAction; the per-item filter.Check, the
// CanSearchDalConnections access guard and the proc enrichment live here
// (search is in customAccessOps, so the generated method emits no access check).
func (svc *dalConnection) onSearch(ctx context.Context, filter types.DalConnectionFilter, aProps *dalConnectionActionProps) (r types.DalConnectionSet, f types.DalConnectionFilter, err error) {
	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.DalConnection) (bool, error) {
		if !svc.ac.CanReadDalConnection(ctx, res) {
			return false, nil
		}

		return true, nil
	}

	if !svc.ac.CanSearchDalConnections(ctx) {
		return nil, f, DalConnectionErrNotAllowedToSearch()
	}

	if r, f, err = store.SearchDalConnections(ctx, svc.store, filter); err != nil {
		return nil, f, err
	}

	svc.proc(ctx, r...)
	return r, f, nil
}

// onCreate is the custom body for the generated Create. The generated method
// owns the action-log scaffold + recordAction; the missing-name guard (which
// must run before the access check), the CanCreateDalConnection access guard
// (create is in customAccessOps), type validation, the dal manager replace
// side-effect and proc enrichment live here.
func (svc *dalConnection) onCreate(ctx context.Context, new *types.DalConnection) (err error) {
	if new.Meta.Name == "" {
		return DalConnectionErrMissingName()
	}

	if !svc.ac.CanCreateDalConnection(ctx) {
		return DalConnectionErrNotAllowedToCreate()
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

	if err = dalConnectionReplace(ctx, svc.store.ToDalConn(), svc.dal, new); err != nil {
		return err
	}
	svc.proc(ctx, new)
	return
}

// onUpdate is the custom body for the generated Update. The generated method
// owns the action-log scaffold + recordAction; the bespoke load/stale/access
// flow, the primary-connection DAL-config handling (with the
// CanManageDalConfigOnDalConnection check), the dal manager replace side-effect
// and proc enrichment live here.
func (svc *dalConnection) onUpdate(ctx context.Context, upd *types.DalConnection, aProps *dalConnectionActionProps) (q *types.DalConnection, err error) {
	var (
		old *types.DalConnection
	)

	err = func() (err error) {
		if upd.Meta.Name == "" {
			return DalConnectionErrMissingName()
		}

		if old, err = loadDalConnection(ctx, svc.store, upd.ID); err != nil {
			return DalConnectionErrNotFound(aProps)
		}

		if !svc.ac.CanUpdateDalConnection(ctx, old) {
			return DalConnectionErrNotAllowedToUpdate(aProps)
		}

		// Test if stale (update has an older version of data)
		if isStale(upd.UpdatedAt, old.UpdatedAt, old.CreatedAt) {
			return DalConnectionErrStaleData()
		}

		upd.UpdatedAt = now()
		upd.CreatedAt = old.CreatedAt
		upd.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()

		// validate
		{
			if old.Type == types.DalPrimaryConnectionResourceType {
				// when primary connection is updated,
				// ignore connection & DAL config changes
				//
				// see Test_dal_connection_update_primary
				// for more details
				upd.Config.DAL = old.Config.DAL
			} else if upd.Config.DAL == nil {
				upd.Config.DAL = old.Config.DAL
			} else if !svc.ac.CanManageDalConfigOnDalConnection(ctx, old) {
				return DalConnectionErrNotAllowedToUpdate()
			}
		}

		if err = store.UpdateDalConnection(ctx, svc.store, upd); err != nil {
			return
		}

		q = upd

		return dalConnectionReplace(ctx, svc.store.ToDalConn(), svc.dal, upd)
	}()

	if q != nil {
		svc.proc(ctx, q)
	}

	return q, err
}

// onDelete is the custom body for the generated DeleteByID. The generated method
// owns the action-log scaffold + recordAction; the primary-connection guard, the
// CanDeleteDalConnection access check, the soft-delete write and the dal manager
// remove side-effect live here.
func (svc *dalConnection) onDelete(ctx context.Context, ID uint64, aProps *dalConnectionActionProps) (err error) {
	var (
		c *types.DalConnection
	)

	if c, err = loadDalConnection(ctx, svc.store, ID); err != nil {
		return
	}

	if c.Type == types.DalPrimaryConnectionResourceType {
		return fmt.Errorf("not allowed to delete primary connections")
	}

	if !svc.ac.CanDeleteDalConnection(ctx, c) {
		return DalConnectionErrNotAllowedToDelete(aProps)
	}

	aProps.setConnection(c)

	c.DeletedAt = now()
	c.DeletedBy = a.GetIdentityFromContext(ctx).Identity()

	if err = store.UpdateDalConnection(ctx, svc.store, c); err != nil {
		return
	}

	return dalConnectionRemove(ctx, svc.dal, c)
}

func (svc *dalConnection) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		cProps = &dalConnectionActionProps{}
		c      *types.DalConnection
	)

	err = func() (err error) {
		if c, err = loadDalConnection(ctx, svc.store, ID); err != nil {
			return
		}

		if !svc.ac.CanDeleteDalConnection(ctx, c) {
			return DalConnectionErrNotAllowedToUndelete(cProps)
		}

		cProps.setConnection(c)

		c.DeletedAt = nil
		c.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()

		if err = store.UpdateDalConnection(ctx, svc.store, c); err != nil {
			return
		}

		// We're creating it here since it was removed on delete
		// primary connection can't be deleted we're just using nil here.
		return dalConnectionReplace(ctx, nil, svc.dal, c)
	}()

	return svc.recordAction(ctx, cProps, DalConnectionActionDelete, err)
}

func (svc *dalConnection) ReloadConnections(ctx context.Context) (err error) {
	return dalConnectionReload(ctx, svc.store, svc.dal)
}

// proc is a helper function that processes the given connection set
// before connections are returned to the caller.
func (svc *dalConnection) proc(ctx context.Context, connections ...*types.DalConnection) {
	for _, c := range connections {
		svc.procPrimaryConnection(c)
		svc.procDal(ctx, c)
		svc.procLocale(c)
	}
}

func (svc *dalConnection) procPrimaryConnection(c *types.DalConnection) {
	if c.Type == types.DalPrimaryConnectionResourceType {
		if c.Config.DAL == nil {
			c.Config.DAL = &types.DalConnectionConfigDAL{}
		}

		return
	}
}

func (svc *dalConnection) procDal(ctx context.Context, c *types.DalConnection) {
	if !svc.ac.CanManageDalConfigOnDalConnection(ctx, c) {
		c.Config.DAL = nil
		return
	}

	c.Issues = svc.dal.SearchConnectionIssues(c.ID)
	if len(c.Issues) == 0 {
		c.Issues = nil
	}
}

func (svc *dalConnection) procLocale(c *types.DalConnection) {
	// @todo...
}

func loadDalConnection(ctx context.Context, s store.DalConnections, ID uint64) (res *types.DalConnection, err error) {
	if ID == 0 {
		return nil, DalConnectionErrInvalidID()
	}

	if res, err = store.LookupDalConnectionByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, DalConnectionErrNotFound()
	}

	return
}

func dalConnectionReload(ctx context.Context, s store.Storer, dcm dalConnManager) (err error) {
	// Get all available connections
	cc, _, err := store.SearchDalConnections(ctx, s, types.DalConnectionFilter{})
	if err != nil {
		return
	}

	return dalConnectionReplace(ctx, s.ToDalConn(), dcm, cc...)
}

// Replaces all given connections
func dalConnectionReplace(ctx context.Context, primary dal.Connection, dcm dalConnManager, cc ...*types.DalConnection) (err error) {
	var (
		cw        *dal.ConnectionWrap
		isPrimary bool
	)

	for _, c := range cc {
		isPrimary = c.Type == types.DalPrimaryConnectionResourceType

		if isPrimary {
			// reuse primary connection
			cw, err = MakeDalConnection(c, primary)
		} else {
			cw, err = MakeDalConnection(c, nil)
		}

		if err != nil {
			return
		}

		if err = dcm.ReplaceConnection(ctx, cw, isPrimary); err != nil {
			return
		}
	}

	return
}

// MakeDalConnection converts types.DalConnection to dal.ConnectionWrap and returns it.
func MakeDalConnection(c *types.DalConnection, existing dal.Connection) (cw *dal.ConnectionWrap, err error) {
	var (
		connConfig = dal.ConnectionConfig{
			SensitivityLevelID: c.Config.Privacy.SensitivityLevelID,
			Label:              c.Handle,
		}
		connParams dal.ConnectionParams
	)

	if c.Config.DAL != nil {
		connConfig.ModelIdent = c.Config.DAL.ModelIdent

		// Copy params and add connection ID
		params := make(map[string]any)
		for k, v := range c.Config.DAL.Params {
			params[k] = v
		}
		params["connectionID"] = c.ID

		connParams = dal.ConnectionParams{
			Type:   c.Config.DAL.Type,
			Params: params,
		}

		if checks := len(c.Config.DAL.ModelIdentCheck); checks > 0 {
			connConfig.ModelIdentCheck = make([]*regexp.Regexp, checks)
			for i, m := range c.Config.DAL.ModelIdentCheck {
				if connConfig.ModelIdentCheck[i], err = regexp.Compile(m); err != nil {
					return nil, fmt.Errorf("could not prepare connection model ident check for %q: %w", c.Handle, err)
				}
			}
		}
	}

	cw = dal.MakeConnection(
		c.ID,
		existing,
		connParams,
		connConfig,
	)

	return
}

// Removes a connection from DAL service
func dalConnectionRemove(ctx context.Context, dcm dalConnManager, cc ...*types.DalConnection) (err error) {
	for _, c := range cc {
		if err = dcm.RemoveConnection(ctx, c.ID); err != nil {
			return err
		}
	}

	return
}
