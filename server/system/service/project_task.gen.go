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

type projectTaskAccessController interface {
	CanCreateProjectTask(context.Context) bool
	CanSearchProjectTasks(context.Context) bool
	CanReadProjectTask(context.Context, *types.ProjectTask) bool
	CanUpdateProjectTask(context.Context, *types.ProjectTask) bool
	CanDeleteProjectTask(context.Context, *types.ProjectTask) bool
}

type projectTask struct {
	actionlog actionlog.Recorder
	store     store.Storer
	ac        projectTaskAccessController
}

type projectTaskServices struct {
	scope scope.Scope
	caps  scope.Capabilities
}

func (svc *projectTask) FindByID(ctx context.Context, ID uint64) (res *types.ProjectTask, err error) {
	var (
		aProps = &projectTaskActionProps{projectTask: &types.ProjectTask{ID: ID}}
	)

	err = func() error {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}
		if res, err = loadProjectTask(ctx, svc.store, ID); err != nil {
			return ProjectTaskErrInvalidID().Wrap(err)
		}

		aProps.setProjectTask(res)

		if err = svc.guard(ctx, res); err != nil {
			return err
		}

		if !svc.ac.CanReadProjectTask(ctx, res) {
			return ProjectTaskErrNotAllowedToRead()
		}

		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ProjectTaskActionLookup, err)
}

func (svc *projectTask) Search(ctx context.Context, filter types.ProjectTaskFilter) (set types.ProjectTaskSet, f types.ProjectTaskFilter, err error) {
	var (
		aProps = &projectTaskActionProps{filter: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.ProjectTask) (bool, error) {
		if !svc.ac.CanReadProjectTask(ctx, res) {
			return false, nil
		}

		return true, nil
	}

	err = func() error {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}
		if !svc.ac.CanSearchProjectTasks(ctx) {
			return ProjectTaskErrNotAllowedToSearch()
		}

		if set, f, err = store.SearchProjectTasks(ctx, svc.store, filter); err != nil {
			return err
		}

		return nil
	}()

	return set, f, svc.recordAction(ctx, aProps, ProjectTaskActionSearch, err)
}

func (svc *projectTask) Create(ctx context.Context, new *types.ProjectTask) (res *types.ProjectTask, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &projectTaskActionProps{projectTask: new, new: new}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}
		if !svc.ac.CanCreateProjectTask(ctx) {
			return ProjectTaskErrNotAllowedToCreate()
		}

		if err = svc.beforeCreate(ctx, new); err != nil {
			return err
		}

		new.ID = nextID()
		new.CreatedAt = *now()

		if err = store.CreateProjectTask(ctx, svc.store, new); err != nil {
			return
		}

		res = new
		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ProjectTaskActionCreate, err)
}

func (svc *projectTask) Update(ctx context.Context, upd *types.ProjectTask) (res *types.ProjectTask, err error) {
	var (
		aProps = &projectTaskActionProps{update: upd}
		old    *types.ProjectTask
	)
	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}
		if res, err = loadProjectTask(ctx, svc.store, upd.ID); err != nil {
			return
		}

		old = res.Clone()
		aProps.setProjectTask(res)

		if err = svc.guard(ctx, res); err != nil {
			return err
		}

		if !svc.ac.CanUpdateProjectTask(ctx, res) {
			return ProjectTaskErrNotAllowedToUpdate()
		}

		// Test if stale (update has an older version of data)
		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return ProjectTaskErrStaleData()
		}

		if err = svc.beforeUpdate(ctx, upd, res); err != nil {
			return err
		}
		res.RevisionID = upd.RevisionID
		res.Title = upd.Title
		res.Description = upd.Description
		res.TaskName = upd.TaskName
		res.TaskType = upd.TaskType
		res.Status = upd.Status
		res.Severity = upd.Severity
		res.Risk = upd.Risk
		res.Owner = upd.Owner
		res.ChangeOwner = upd.ChangeOwner
		res.DateDue = upd.DateDue
		res.CompletedDate = upd.CompletedDate
		res.UpdatedAt = now()

		if err = store.UpdateProjectTask(ctx, svc.store, res); err != nil {
			return err
		}
		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ProjectTaskActionUpdate, err, old, res)
}

func (svc *projectTask) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &projectTaskActionProps{}
		res    *types.ProjectTask
	)
	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		if res, err = loadProjectTask(ctx, svc.store, ID); err != nil {
			return
		}

		if err = svc.guard(ctx, res); err != nil {
			return err
		}

		aProps.setProjectTask(res)

		if !svc.ac.CanDeleteProjectTask(ctx, res) {
			return ProjectTaskErrNotAllowedToDelete()
		}

		if err = svc.beforeDelete(ctx, res); err != nil {
			return err
		}

		res.DeletedAt = now()
		if err = store.UpdateProjectTask(ctx, svc.store, res); err != nil {
			return
		}
		return nil
	}()

	return svc.recordAction(ctx, aProps, ProjectTaskActionDelete, err)
}

func loadProjectTask(ctx context.Context, s store.ProjectTasks, ID uint64) (res *types.ProjectTask, err error) {
	if ID == 0 {
		return nil, ProjectTaskErrInvalidID()
	}

	if res, err = store.LookupProjectTaskByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, ProjectTaskErrNotFound()
	}

	return
}
func (svc *projectTask) guard(_ context.Context, _ *types.ProjectTask) error { return nil }

func (svc *projectTask) checkScope(ctx context.Context, cap scope.Capability) error {
	if err := scope.RequireTenantMembership(ctx); err != nil {
		return err
	}
	return scope.RequireCapability(ctx, cap)
}
func (svc *projectTask) scopeServices(ctx context.Context) *projectTaskServices {
	return &projectTaskServices{
		scope: scope.GetScopeFromContext(ctx),
		caps:  scope.GetCapabilitiesFromContext(ctx),
	}
}
