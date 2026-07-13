package service

import (
	"context"
	"fmt"
	"regexp"

	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/dal"
	"github.com/crusttech/human/server/pkg/options"

	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

type (
	dalConnectionServices struct {
		dal    dalConnManager
		dbConf options.DBOpt
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
		services: &dalConnectionServices{
			dal:    dal,
			dbConf: dbConf,
		},
	}
}

func (svc *dalConnection) onLookup(ctx context.Context, ID uint64, aProps *dalConnectionActionProps) (*types.DalConnection, error) {
	res, err := loadDalConnection(ctx, svc.store, ID)
	if err != nil {
		return nil, DalConnectionErrInvalidID().Wrap(err)
	}

	aProps.setConnection(res)

	if !svc.ac.CanReadDalConnection(ctx, res) {
		return nil, DalConnectionErrNotAllowedToRead(aProps)
	}

	svc.proc(ctx, res)
	return res, nil
}

func (svc *dalConnection) onSearch(ctx context.Context, filter types.DalConnectionFilter, aProps *dalConnectionActionProps) (types.DalConnectionSet, types.DalConnectionFilter, error) {
	filter.Check = func(res *types.DalConnection) (bool, error) {
		if !svc.ac.CanReadDalConnection(ctx, res) {
			return false, nil
		}
		return true, nil
	}

	if !svc.ac.CanSearchDalConnections(ctx) {
		return nil, filter, DalConnectionErrNotAllowedToSearch()
	}

	set, f, err := store.SearchDalConnections(ctx, svc.store, filter)
	if err != nil {
		return nil, f, err
	}

	svc.proc(ctx, set...)
	return set, f, nil
}

func (svc *dalConnection) onCreate(ctx context.Context, new *types.DalConnection) error {
	if new.Meta.Name == "" {
		return DalConnectionErrMissingName()
	}

	if !svc.ac.CanCreateDalConnection(ctx) {
		return DalConnectionErrNotAllowedToCreate()
	}

	new.ID = nextID()
	new.CreatedAt = *now()
	new.CreatedBy = a.GetIdentityFromContext(ctx).Identity()

	if new.Type != types.DalConnectionResourceType {
		return fmt.Errorf("cannot create connection: unsupported connection type %s", new.Type)
	}

	if err := store.CreateDalConnection(ctx, svc.store, new); err != nil {
		return err
	}

	if err := dalConnectionReplace(ctx, svc.store.ToDalConn(), svc.services.dal, new); err != nil {
		return err
	}

	svc.proc(ctx, new)
	return nil
}

func (svc *dalConnection) onUpdate(ctx context.Context, s store.Storer, upd *types.DalConnection, res *types.DalConnection, aProps *dalConnectionActionProps, before func() error, after func() error) error {
	if upd.Meta.Name == "" {
		return DalConnectionErrMissingName()
	}

	if !svc.ac.CanUpdateDalConnection(ctx, res) {
		return DalConnectionErrNotAllowedToUpdate(aProps)
	}

	if res.Type == types.DalPrimaryConnectionResourceType {
		upd.Config.DAL = res.Config.DAL
	} else if upd.Config.DAL == nil {
		upd.Config.DAL = res.Config.DAL
	} else if !svc.ac.CanManageDalConfigOnDalConnection(ctx, res) {
		return DalConnectionErrNotAllowedToUpdate()
	}

	upd.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()

	if err := store.UpdateDalConnection(ctx, s, upd); err != nil {
		return err
	}

	svc.proc(ctx, upd)
	return dalConnectionReplace(ctx, svc.store.ToDalConn(), svc.services.dal, upd)
}

func (svc *dalConnection) onDelete(ctx context.Context, s store.Storer, res *types.DalConnection, aProps *dalConnectionActionProps) error {
	if res.Type == types.DalPrimaryConnectionResourceType {
		return fmt.Errorf("not allowed to delete primary connections")
	}

	if !svc.ac.CanDeleteDalConnection(ctx, res) {
		return DalConnectionErrNotAllowedToDelete(aProps)
	}

	res.DeletedAt = now()
	res.DeletedBy = a.GetIdentityFromContext(ctx).Identity()

	if err := store.UpdateDalConnection(ctx, s, res); err != nil {
		return err
	}

	return dalConnectionRemove(ctx, svc.services.dal, res)
}

func (svc *dalConnection) onReloadConnections(ctx context.Context, aProps *dalConnectionActionProps) error {
	return dalConnectionReload(ctx, svc.store, svc.services.dal)
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

		// primary connection can't be deleted; passing nil for primary conn here
		return dalConnectionReplace(ctx, nil, svc.services.dal, c)
	}()

	return svc.recordAction(ctx, cProps, DalConnectionActionDelete, err)
}

// proc processes the given connections before they are returned to the caller.
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
	}
}

func (svc *dalConnection) procDal(ctx context.Context, c *types.DalConnection) {
	if !svc.ac.CanManageDalConfigOnDalConnection(ctx, c) {
		c.Config.DAL = nil
		return
	}

	c.Issues = svc.services.dal.SearchConnectionIssues(c.ID)
	if len(c.Issues) == 0 {
		c.Issues = nil
	}
}

func (svc *dalConnection) procLocale(c *types.DalConnection) {
	// @todo...
}

func dalConnectionReload(ctx context.Context, s store.Storer, dcm dalConnManager) (err error) {
	cc, _, err := store.SearchDalConnections(ctx, s, types.DalConnectionFilter{})
	if err != nil {
		return
	}

	return dalConnectionReplace(ctx, s.ToDalConn(), dcm, cc...)
}

// dalConnectionReplace replaces all given connections in the DAL service.
func dalConnectionReplace(ctx context.Context, primary dal.Connection, dcm dalConnManager, cc ...*types.DalConnection) (err error) {
	var (
		cw        *dal.ConnectionWrap
		isPrimary bool
	)

	for _, c := range cc {
		isPrimary = c.Type == types.DalPrimaryConnectionResourceType

		if isPrimary {
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

// MakeDalConnection converts types.DalConnection to dal.ConnectionWrap.
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

		connParams = dal.ConnectionParams{
			Type:   c.Config.DAL.Type,
			Params: c.Config.DAL.Params,
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

// dalConnectionRemove removes connections from the DAL service.
func dalConnectionRemove(ctx context.Context, dcm dalConnManager, cc ...*types.DalConnection) (err error) {
	for _, c := range cc {
		if err = dcm.RemoveConnection(ctx, c.ID); err != nil {
			return err
		}
	}

	return
}
