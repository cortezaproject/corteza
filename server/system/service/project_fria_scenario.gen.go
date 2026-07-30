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

type projectFriaScenarioAccessController interface {
	CanCreateProjectFriaScenario(context.Context) bool
	CanSearchProjectFriaScenarios(context.Context) bool
	CanReadProjectFriaScenario(context.Context, *types.ProjectFriaScenario) bool
	CanUpdateProjectFriaScenario(context.Context, *types.ProjectFriaScenario) bool
	CanDeleteProjectFriaScenario(context.Context, *types.ProjectFriaScenario) bool
}

type projectFriaScenario struct {
	actionlog actionlog.Recorder
	store     store.Storer
	ac        projectFriaScenarioAccessController
}

type projectFriaScenarioServices struct {
	scope scope.Scope
	caps  scope.Capabilities
}

func (svc *projectFriaScenario) FindByID(ctx context.Context, ID uint64) (res *types.ProjectFriaScenario, err error) {
	var (
		aProps = &projectFriaScenarioActionProps{projectFriaScenario: &types.ProjectFriaScenario{ID: ID}}
	)

	err = func() error {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}
		if res, err = loadProjectFriaScenario(ctx, svc.store, ID); err != nil {
			return ProjectFriaScenarioErrInvalidID().Wrap(err)
		}

		aProps.setProjectFriaScenario(res)

		if err = svc.guard(ctx, res); err != nil {
			return err
		}

		if !svc.ac.CanReadProjectFriaScenario(ctx, res) {
			return ProjectFriaScenarioErrNotAllowedToRead()
		}

		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ProjectFriaScenarioActionLookup, err)
}

func (svc *projectFriaScenario) Search(ctx context.Context, filter types.ProjectFriaScenarioFilter) (set types.ProjectFriaScenarioSet, f types.ProjectFriaScenarioFilter, err error) {
	var (
		aProps = &projectFriaScenarioActionProps{search: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.ProjectFriaScenario) (bool, error) {
		if !svc.ac.CanReadProjectFriaScenario(ctx, res) {
			return false, nil
		}

		return true, nil
	}

	err = func() error {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}
		if !svc.ac.CanSearchProjectFriaScenarios(ctx) {
			return ProjectFriaScenarioErrNotAllowedToSearch()
		}

		if set, f, err = store.SearchProjectFriaScenarios(ctx, svc.store, filter); err != nil {
			return err
		}

		return nil
	}()

	return set, f, svc.recordAction(ctx, aProps, ProjectFriaScenarioActionSearch, err)
}

func (svc *projectFriaScenario) Create(ctx context.Context, new *types.ProjectFriaScenario) (res *types.ProjectFriaScenario, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &projectFriaScenarioActionProps{projectFriaScenario: new, new: new}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}
		if !svc.ac.CanCreateProjectFriaScenario(ctx) {
			return ProjectFriaScenarioErrNotAllowedToCreate()
		}

		if err = svc.beforeCreate(ctx, new); err != nil {
			return err
		}

		new.ID = nextID()
		new.CreatedAt = *now()

		if err = store.CreateProjectFriaScenario(ctx, svc.store, new); err != nil {
			return
		}

		res = new
		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ProjectFriaScenarioActionCreate, err)
}

func (svc *projectFriaScenario) Update(ctx context.Context, upd *types.ProjectFriaScenario) (res *types.ProjectFriaScenario, err error) {
	var (
		aProps = &projectFriaScenarioActionProps{update: upd}
		old    *types.ProjectFriaScenario
	)
	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}
		if res, err = loadProjectFriaScenario(ctx, svc.store, upd.ID); err != nil {
			return
		}

		old = res.Clone()
		aProps.setProjectFriaScenario(res)

		if err = svc.guard(ctx, res); err != nil {
			return err
		}

		if !svc.ac.CanUpdateProjectFriaScenario(ctx, res) {
			return ProjectFriaScenarioErrNotAllowedToUpdate()
		}

		// Test if stale (update has an older version of data)
		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return ProjectFriaScenarioErrStaleData()
		}

		if err = svc.beforeUpdate(ctx, upd, res); err != nil {
			return err
		}
		res.AiSystemID = upd.AiSystemID
		res.Title = upd.Title
		res.Severity = upd.Severity
		res.UpdatedBy = upd.UpdatedBy
		res.UpdatedAt = now()

		if err = store.UpdateProjectFriaScenario(ctx, svc.store, res); err != nil {
			return err
		}
		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ProjectFriaScenarioActionUpdate, err, old, res)
}

func (svc *projectFriaScenario) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &projectFriaScenarioActionProps{}
		res    *types.ProjectFriaScenario
	)
	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		if res, err = loadProjectFriaScenario(ctx, svc.store, ID); err != nil {
			return
		}

		if err = svc.guard(ctx, res); err != nil {
			return err
		}

		aProps.setProjectFriaScenario(res)

		if !svc.ac.CanDeleteProjectFriaScenario(ctx, res) {
			return ProjectFriaScenarioErrNotAllowedToDelete()
		}

		if err = svc.beforeDelete(ctx, res); err != nil {
			return err
		}

		res.DeletedAt = now()
		if err = store.UpdateProjectFriaScenario(ctx, svc.store, res); err != nil {
			return
		}
		return nil
	}()

	return svc.recordAction(ctx, aProps, ProjectFriaScenarioActionDelete, err)
}

func loadProjectFriaScenario(ctx context.Context, s store.ProjectFriaScenarios, ID uint64) (res *types.ProjectFriaScenario, err error) {
	if ID == 0 {
		return nil, ProjectFriaScenarioErrInvalidID()
	}

	if res, err = store.LookupProjectFriaScenarioByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, ProjectFriaScenarioErrNotFound()
	}

	return
}
func (svc *projectFriaScenario) guard(_ context.Context, _ *types.ProjectFriaScenario) error {
	return nil
}

func (svc *projectFriaScenario) checkScope(ctx context.Context, cap scope.Capability) error {
	if err := scope.RequireTenantMembership(ctx); err != nil {
		return err
	}
	return scope.RequireCapability(ctx, cap)
}
func (svc *projectFriaScenario) scopeServices(ctx context.Context) *projectFriaScenarioServices {
	return &projectFriaScenarioServices{
		scope: scope.GetScopeFromContext(ctx),
		caps:  scope.GetCapabilitiesFromContext(ctx),
	}
}
