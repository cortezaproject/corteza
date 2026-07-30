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
	"github.com/crusttech/human/server/pkg/handle"
	"github.com/crusttech/human/server/pkg/scope"
	"github.com/crusttech/human/server/store"
	types "github.com/crusttech/human/server/system/types"
)

type projectAiSystem struct {
	actionlog actionlog.Recorder
	store     store.Storer
	ac        projectAiSystemAccessController
}

type projectAiSystemServices struct {
	scope scope.Scope
	caps  scope.Capabilities
}

func (svc *projectAiSystem) FindByID(ctx context.Context, ID uint64) (res *types.ProjectAiSystem, err error) {
	var (
		aProps = &projectAiSystemActionProps{projectAiSystem: &types.ProjectAiSystem{ID: ID}}
	)

	err = func() error {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}
		if res, err = loadProjectAiSystem(ctx, svc.store, ID); err != nil {
			return ProjectAiSystemErrInvalidID().Wrap(err)
		}

		aProps.setProjectAiSystem(res)

		if err = svc.guard(ctx, res); err != nil {
			return err
		}

		if !svc.ac.CanReadProjectAiSystem(ctx, res) {
			return ProjectAiSystemErrNotAllowedToRead()
		}

		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ProjectAiSystemActionLookup, err)
}

func (svc *projectAiSystem) Search(ctx context.Context, filter types.ProjectAiSystemFilter) (set types.ProjectAiSystemSet, f types.ProjectAiSystemFilter, err error) {
	var (
		aProps = &projectAiSystemActionProps{search: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.ProjectAiSystem) (bool, error) {
		if !svc.ac.CanReadProjectAiSystem(ctx, res) {
			return false, nil
		}

		return true, nil
	}

	err = func() error {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}
		if !svc.ac.CanSearchProjectAiSystems(ctx) {
			return ProjectAiSystemErrNotAllowedToSearch()
		}

		if set, f, err = store.SearchProjectAiSystems(ctx, svc.store, filter); err != nil {
			return err
		}

		return nil
	}()

	return set, f, svc.recordAction(ctx, aProps, ProjectAiSystemActionSearch, err)
}

func (svc *projectAiSystem) Create(ctx context.Context, new *types.ProjectAiSystem) (res *types.ProjectAiSystem, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &projectAiSystemActionProps{projectAiSystem: new, new: new}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}
		if !svc.ac.CanCreateProjectAiSystem(ctx) {
			return ProjectAiSystemErrNotAllowedToCreate()
		}

		if err = svc.beforeCreate(ctx, new); err != nil {
			return err
		}

		new.ID = nextID()
		new.CreatedAt = *now()

		if err = store.CreateProjectAiSystem(ctx, svc.store, new); err != nil {
			return
		}

		res = new
		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ProjectAiSystemActionCreate, err)
}

func (svc *projectAiSystem) Update(ctx context.Context, upd *types.ProjectAiSystem) (res *types.ProjectAiSystem, err error) {
	var (
		aProps = &projectAiSystemActionProps{update: upd}
		old    *types.ProjectAiSystem
	)
	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}
		if res, err = loadProjectAiSystem(ctx, svc.store, upd.ID); err != nil {
			return
		}

		old = res.Clone()
		aProps.setProjectAiSystem(res)

		if err = svc.guard(ctx, res); err != nil {
			return err
		}

		if upd.Handle != res.Handle && !handle.IsValid(upd.Handle) {
			return ProjectAiSystemErrInvalidHandle()
		}

		if !svc.ac.CanUpdateProjectAiSystem(ctx, res) {
			return ProjectAiSystemErrNotAllowedToUpdate()
		}

		// Test if stale (update has an older version of data)
		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return ProjectAiSystemErrStaleData()
		}

		if err = svc.beforeUpdate(ctx, upd, res); err != nil {
			return err
		}
		res.RiskClass = upd.RiskClass
		res.UpdatedAt = now()

		if err = store.UpdateProjectAiSystem(ctx, svc.store, res); err != nil {
			return err
		}
		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ProjectAiSystemActionUpdate, err, old, res)
}

func (svc *projectAiSystem) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &projectAiSystemActionProps{}
		res    *types.ProjectAiSystem
	)
	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		if res, err = loadProjectAiSystem(ctx, svc.store, ID); err != nil {
			return
		}

		if err = svc.guard(ctx, res); err != nil {
			return err
		}

		aProps.setProjectAiSystem(res)

		if !svc.ac.CanDeleteProjectAiSystem(ctx, res) {
			return ProjectAiSystemErrNotAllowedToDelete()
		}

		res.DeletedAt = now()
		if err = store.UpdateProjectAiSystem(ctx, svc.store, res); err != nil {
			return
		}
		return nil
	}()

	return svc.recordAction(ctx, aProps, ProjectAiSystemActionDelete, err)
}

func loadProjectAiSystem(ctx context.Context, s store.ProjectAiSystems, ID uint64) (res *types.ProjectAiSystem, err error) {
	if ID == 0 {
		return nil, ProjectAiSystemErrInvalidID()
	}

	if res, err = store.LookupProjectAiSystemByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, ProjectAiSystemErrNotFound()
	}

	return
}
func (svc *projectAiSystem) guard(_ context.Context, _ *types.ProjectAiSystem) error { return nil }

func (svc *projectAiSystem) checkScope(ctx context.Context, cap scope.Capability) error {
	if err := scope.RequireTenantMembership(ctx); err != nil {
		return err
	}
	return scope.RequireCapability(ctx, cap)
}
func (svc *projectAiSystem) scopeServices(ctx context.Context) *projectAiSystemServices {
	return &projectAiSystemServices{
		scope: scope.GetScopeFromContext(ctx),
		caps:  scope.GetCapabilitiesFromContext(ctx),
	}
}

func (svc *projectAiSystem) MemberList(ctx context.Context, projectAiSystemID uint64) (set types.ProjectAiSystemEntrySet, err error) {
	var (
		aProps = &projectAiSystemActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}

		set, err = svc.onMemberList(ctx, aProps, projectAiSystemID)
		return err
	}()

	return set, svc.recordAction(ctx, aProps, ProjectAiSystemActionMemberList, err)
}

func (svc *projectAiSystem) MemberAdd(ctx context.Context, projectAiSystemID uint64, resourceRef string) (err error) {
	var (
		aProps = &projectAiSystemActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		err = svc.onMemberAdd(ctx, aProps, projectAiSystemID, resourceRef)
		return err
	}()

	return svc.recordAction(ctx, aProps, ProjectAiSystemActionMemberAdd, err)
}

func (svc *projectAiSystem) MemberRemove(ctx context.Context, projectAiSystemID uint64, resourceRef string) (err error) {
	var (
		aProps = &projectAiSystemActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		err = svc.onMemberRemove(ctx, aProps, projectAiSystemID, resourceRef)
		return err
	}()

	return svc.recordAction(ctx, aProps, ProjectAiSystemActionMemberRemove, err)
}
