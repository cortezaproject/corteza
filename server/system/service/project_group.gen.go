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

type projectGroup struct {
	actionlog actionlog.Recorder
	store     store.Storer
	ac        projectGroupAccessController
}

func ProjectGroup() *projectGroup {
	return &projectGroup{
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		ac:        DefaultAccessControl,
	}
}

type projectGroupServices struct {
	scope scope.Scope
	caps  scope.Capabilities
}

func (svc *projectGroup) FindByID(ctx context.Context, ID uint64) (res *types.ProjectGroup, err error) {
	var (
		aProps = &projectGroupActionProps{projectGroup: &types.ProjectGroup{ID: ID}}
	)

	err = func() error {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}
		if res, err = loadProjectGroup(ctx, svc.store, ID); err != nil {
			return ProjectGroupErrInvalidID().Wrap(err)
		}

		aProps.setProjectGroup(res)

		if !svc.ac.CanReadProjectGroup(ctx, res) {
			return ProjectGroupErrNotAllowedToRead()
		}

		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ProjectGroupActionLookup, err)
}

func (svc *projectGroup) Search(ctx context.Context, filter types.ProjectGroupFilter) (set types.ProjectGroupSet, f types.ProjectGroupFilter, err error) {
	var (
		aProps = &projectGroupActionProps{search: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.ProjectGroup) (bool, error) {
		if !svc.ac.CanReadProjectGroup(ctx, res) {
			return false, nil
		}

		return true, nil
	}

	err = func() error {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}
		if !svc.ac.CanSearchProjectGroups(ctx) {
			return ProjectGroupErrNotAllowedToSearch()
		}

		if set, f, err = store.SearchProjectGroups(ctx, svc.store, filter); err != nil {
			return err
		}

		return nil
	}()

	return set, f, svc.recordAction(ctx, aProps, ProjectGroupActionSearch, err)
}

func (svc *projectGroup) Create(ctx context.Context, new *types.ProjectGroup) (res *types.ProjectGroup, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &projectGroupActionProps{projectGroup: new, new: new}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}
		if !svc.ac.CanCreateProjectGroup(ctx) {
			return ProjectGroupErrNotAllowedToCreate()
		}

		if err = svc.beforeCreate(ctx, new); err != nil {
			return err
		}

		new.ID = nextID()
		new.CreatedAt = *now()

		if err = store.CreateProjectGroup(ctx, svc.store, new); err != nil {
			return
		}

		res = new
		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ProjectGroupActionCreate, err)
}

func (svc *projectGroup) Update(ctx context.Context, upd *types.ProjectGroup) (res *types.ProjectGroup, err error) {
	var (
		aProps = &projectGroupActionProps{update: upd}
		old    *types.ProjectGroup
	)
	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}
		if res, err = loadProjectGroup(ctx, svc.store, upd.ID); err != nil {
			return
		}

		old = res.Clone()
		aProps.setProjectGroup(res)

		if upd.Handle != res.Handle && !handle.IsValid(upd.Handle) {
			return ProjectGroupErrInvalidHandle()
		}

		if !svc.ac.CanUpdateProjectGroup(ctx, res) {
			return ProjectGroupErrNotAllowedToUpdate()
		}

		// Test if stale (update has an older version of data)
		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return ProjectGroupErrStaleData()
		}

		if err = svc.beforeUpdate(ctx, upd, res); err != nil {
			return err
		}
		res.UpdatedAt = now()

		if err = store.UpdateProjectGroup(ctx, svc.store, res); err != nil {
			return err
		}
		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ProjectGroupActionUpdate, err, old, res)
}

func (svc *projectGroup) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &projectGroupActionProps{}
		res    *types.ProjectGroup
	)
	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		if res, err = loadProjectGroup(ctx, svc.store, ID); err != nil {
			return
		}

		aProps.setProjectGroup(res)

		if !svc.ac.CanDeleteProjectGroup(ctx, res) {
			return ProjectGroupErrNotAllowedToDelete()
		}

		res.DeletedAt = now()
		if err = store.UpdateProjectGroup(ctx, svc.store, res); err != nil {
			return
		}
		return nil
	}()

	return svc.recordAction(ctx, aProps, ProjectGroupActionDelete, err)
}

func loadProjectGroup(ctx context.Context, s store.ProjectGroups, ID uint64) (res *types.ProjectGroup, err error) {
	if ID == 0 {
		return nil, ProjectGroupErrInvalidID()
	}

	if res, err = store.LookupProjectGroupByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, ProjectGroupErrNotFound()
	}

	return
}

func (svc *projectGroup) checkScope(ctx context.Context, cap scope.Capability) error {
	if err := scope.RequireTenantMembership(ctx); err != nil {
		return err
	}
	return scope.RequireCapability(ctx, cap)
}

func (svc *projectGroup) scopeServices(ctx context.Context) *projectGroupServices {
	return &projectGroupServices{
		scope: scope.GetScopeFromContext(ctx),
		caps:  scope.GetCapabilitiesFromContext(ctx),
	}
}

func (svc *projectGroup) MemberList(ctx context.Context, projectGroupID uint64) (set types.ProjectGroupEntrySet, err error) {
	var (
		aProps = &projectGroupActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}

		set, err = svc.onMemberList(ctx, aProps, projectGroupID)
		return err
	}()

	return set, svc.recordAction(ctx, aProps, ProjectGroupActionMemberList, err)
}

func (svc *projectGroup) MemberAdd(ctx context.Context, projectGroupID uint64, resourceRef string) (err error) {
	var (
		aProps = &projectGroupActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		err = svc.onMemberAdd(ctx, aProps, projectGroupID, resourceRef)
		return err
	}()

	return svc.recordAction(ctx, aProps, ProjectGroupActionMemberAdd, err)
}

func (svc *projectGroup) MemberRemove(ctx context.Context, projectGroupID uint64, resourceRef string) (err error) {
	var (
		aProps = &projectGroupActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		err = svc.onMemberRemove(ctx, aProps, projectGroupID, resourceRef)
		return err
	}()

	return svc.recordAction(ctx, aProps, ProjectGroupActionMemberRemove, err)
}
