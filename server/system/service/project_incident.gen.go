package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"

	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/scope"
	"github.com/crusttech/human/server/store"
	types "github.com/crusttech/human/server/system/types"
)

type projectIncidentAccessController interface {
	CanCreateProjectIncident(context.Context) bool
	CanSearchProjectIncidents(context.Context) bool
	CanReadProjectIncident(context.Context, *types.ProjectIncident) bool
	CanUpdateProjectIncident(context.Context, *types.ProjectIncident) bool
	CanDeleteProjectIncident(context.Context, *types.ProjectIncident) bool
}

type projectIncident struct {
	actionlog actionlog.Recorder
	store     store.Storer
	ac        projectIncidentAccessController
}

type projectIncidentServices struct {
	scope scope.Scope
	caps  scope.Capabilities
}

func (svc *projectIncident) FindByID(ctx context.Context, ID uint64) (res *types.ProjectIncident, err error) {
	var (
		aProps = &projectIncidentActionProps{projectIncident: &types.ProjectIncident{ID: ID}}
	)

	err = func() error {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}
		if res, err = loadProjectIncident(ctx, svc.store, ID); err != nil {
			return ProjectIncidentErrInvalidID().Wrap(err)
		}

		aProps.setProjectIncident(res)

		if err = svc.guard(ctx, res); err != nil {
			return err
		}

		if !svc.ac.CanReadProjectIncident(ctx, res) {
			return ProjectIncidentErrNotAllowedToRead()
		}

		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ProjectIncidentActionLookup, err)
}

func (svc *projectIncident) Search(ctx context.Context, filter types.ProjectIncidentFilter) (set types.ProjectIncidentSet, f types.ProjectIncidentFilter, err error) {
	var (
		aProps = &projectIncidentActionProps{filter: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.ProjectIncident) (bool, error) {
		if !svc.ac.CanReadProjectIncident(ctx, res) {
			return false, nil
		}

		return true, nil
	}

	err = func() error {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}
		if !svc.ac.CanSearchProjectIncidents(ctx) {
			return ProjectIncidentErrNotAllowedToSearch()
		}

		if set, f, err = store.SearchProjectIncidents(ctx, svc.store, filter); err != nil {
			return err
		}

		return nil
	}()

	return set, f, svc.recordAction(ctx, aProps, ProjectIncidentActionSearch, err)
}

func (svc *projectIncident) Create(ctx context.Context, new *types.ProjectIncident) (res *types.ProjectIncident, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &projectIncidentActionProps{projectIncident: new, new: new}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}
		if !svc.ac.CanCreateProjectIncident(ctx) {
			return ProjectIncidentErrNotAllowedToCreate()
		}

		if err = svc.beforeCreate(ctx, new); err != nil {
			return err
		}

		new.ID = nextID()
		new.CreatedAt = *now()

		if err = store.CreateProjectIncident(ctx, svc.store, new); err != nil {
			return
		}

		res = new
		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ProjectIncidentActionCreate, err)
}

func (svc *projectIncident) Update(ctx context.Context, upd *types.ProjectIncident) (res *types.ProjectIncident, err error) {
	var (
		aProps = &projectIncidentActionProps{update: upd}
		old    *types.ProjectIncident
	)
	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}
		if res, err = loadProjectIncident(ctx, svc.store, upd.ID); err != nil {
			return
		}

		old = res.Clone()
		aProps.setProjectIncident(res)

		if err = svc.guard(ctx, res); err != nil {
			return err
		}

		if !svc.ac.CanUpdateProjectIncident(ctx, res) {
			return ProjectIncidentErrNotAllowedToUpdate()
		}

		// Test if stale (update has an older version of data)
		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return ProjectIncidentErrStaleData()
		}

		if err = svc.beforeUpdate(ctx, upd, res); err != nil {
			return err
		}
		res.Title = upd.Title
		res.Description = upd.Description
		res.IncidentType = upd.IncidentType
		res.GroupSystem = upd.GroupSystem
		res.Status = upd.Status
		res.Severity = upd.Severity
		res.Risk = upd.Risk
		res.IssueOwner = upd.IssueOwner
		res.ChangeOwner = upd.ChangeOwner
		res.ChangeApprovedBy = upd.ChangeApprovedBy
		res.RiskIssue = upd.RiskIssue
		res.ChangeRequired = upd.ChangeRequired
		res.RiskChange = upd.RiskChange
		res.Backlog = upd.Backlog
		res.DateDue = upd.DateDue
		res.CompletedDate = upd.CompletedDate
		res.UpdatedAt = now()

		if err = store.UpdateProjectIncident(ctx, svc.store, res); err != nil {
			return err
		}
		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ProjectIncidentActionUpdate, err, old, res)
}

func (svc *projectIncident) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &projectIncidentActionProps{}
		res    *types.ProjectIncident
	)
	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		if res, err = loadProjectIncident(ctx, svc.store, ID); err != nil {
			return
		}

		if err = svc.guard(ctx, res); err != nil {
			return err
		}

		aProps.setProjectIncident(res)

		if !svc.ac.CanDeleteProjectIncident(ctx, res) {
			return ProjectIncidentErrNotAllowedToDelete()
		}

		if err = svc.beforeDelete(ctx, res); err != nil {
			return err
		}

		res.DeletedAt = now()
		if err = store.UpdateProjectIncident(ctx, svc.store, res); err != nil {
			return
		}
		return nil
	}()

	return svc.recordAction(ctx, aProps, ProjectIncidentActionDelete, err)
}

func loadProjectIncident(ctx context.Context, s store.ProjectIncidents, ID uint64) (res *types.ProjectIncident, err error) {
	if ID == 0 {
		return nil, ProjectIncidentErrInvalidID()
	}

	if res, err = store.LookupProjectIncidentByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, ProjectIncidentErrNotFound()
	}

	return
}
func (svc *projectIncident) guard(_ context.Context, _ *types.ProjectIncident) error { return nil }

func (svc *projectIncident) checkScope(ctx context.Context, cap scope.Capability) error {
	if err := scope.RequireTenantMembership(ctx); err != nil {
		return err
	}
	return scope.RequireCapability(ctx, cap)
}
func (svc *projectIncident) scopeServices(ctx context.Context) *projectIncidentServices {
	return &projectIncidentServices{
		scope: scope.GetScopeFromContext(ctx),
		caps:  scope.GetCapabilitiesFromContext(ctx),
	}
}
