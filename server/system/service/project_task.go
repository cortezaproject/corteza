package service

import (
	"context"

	"github.com/crusttech/human/server/pkg/actionlog"
	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

type (
	projectTask struct {
		actionlog actionlog.Recorder
		store     store.Storer
		ac        projectTaskAccessController
	}

	projectTaskAccessController interface {
		CanSearchProjectTasks(context.Context) bool
		CanCreateProjectTask(context.Context) bool
		CanReadProjectTask(context.Context, *types.ProjectTask) bool
		CanUpdateProjectTask(context.Context, *types.ProjectTask) bool
		CanDeleteProjectTask(context.Context, *types.ProjectTask) bool
	}
)

func ProjectTask() *projectTask {
	return &projectTask{
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
	}
}

func (svc projectTask) FindByID(ctx context.Context, ID uint64) (r *types.ProjectTask, err error) {
	var (
		raProps = &projectTaskActionProps{projectTask: &types.ProjectTask{ID: ID}}
	)

	err = func() error {
		if r, err = loadProjectTask(ctx, svc.store, ID); err != nil {
			return err
		}

		raProps.setProjectTask(r)

		if !svc.ac.CanReadProjectTask(ctx, r) {
			return ProjectTaskErrNotAllowedToRead()
		}

		return nil
	}()

	return r, svc.recordAction(ctx, raProps, ProjectTaskActionLookup, err)
}

func (svc projectTask) Search(ctx context.Context, filter types.ProjectTaskFilter) (rr types.ProjectTaskSet, f types.ProjectTaskFilter, err error) {
	var (
		raProps = &projectTaskActionProps{filter: &filter}
	)

	filter.Check = func(res *types.ProjectTask) (bool, error) {
		if !svc.ac.CanReadProjectTask(ctx, res) {
			return false, nil
		}
		return true, nil
	}

	err = func() error {
		if !svc.ac.CanSearchProjectTasks(ctx) {
			return ProjectTaskErrNotAllowedToSearch()
		}

		if rr, f, err = store.SearchProjectTasks(ctx, svc.store, filter); err != nil {
			return err
		}

		return nil
	}()

	return rr, f, svc.recordAction(ctx, raProps, ProjectTaskActionSearch, err)
}

func (svc projectTask) Create(ctx context.Context, new *types.ProjectTask) (r *types.ProjectTask, err error) {
	var (
		raProps = &projectTaskActionProps{projectTask: new, new: new}
	)

	err = func() (err error) {
		if !svc.ac.CanCreateProjectTask(ctx) {
			return ProjectTaskErrNotAllowedToCreate()
		}

		if len(new.Title) == 0 {
			new.Title = firstNonEmpty(new.Description, "(untitled)")
		}
		if len(new.Status) == 0 {
			new.Status = "Open"
		}

		new.ID = nextID()
		new.CreatedAt = *now()
		new.CreatedBy = a.GetIdentityFromContext(ctx).Identity()

		if err = store.CreateProjectTask(ctx, svc.store, new); err != nil {
			return
		}

		r = new
		return nil
	}()

	return r, svc.recordAction(ctx, raProps, ProjectTaskActionCreate, err)
}

func (svc projectTask) Update(ctx context.Context, upd *types.ProjectTask) (r *types.ProjectTask, err error) {
	var (
		raProps = &projectTaskActionProps{projectTask: upd, update: upd}
	)

	err = func() (err error) {
		if r, err = loadProjectTask(ctx, svc.store, upd.ID); err != nil {
			return
		}

		raProps.setProjectTask(r)

		if !svc.ac.CanUpdateProjectTask(ctx, r) {
			return ProjectTaskErrNotAllowedToUpdate()
		}

		r.Title = upd.Title
		r.Description = upd.Description
		r.TaskName = upd.TaskName
		r.TaskType = upd.TaskType
		r.Status = upd.Status
		r.Severity = upd.Severity
		r.Risk = upd.Risk
		r.Owner = upd.Owner
		r.ChangeOwner = upd.ChangeOwner
		r.Backlog = upd.Backlog
		r.DateDue = upd.DateDue
		r.CompletedDate = upd.CompletedDate
		r.UpdatedAt = now()
		r.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()

		if err = store.UpdateProjectTask(ctx, svc.store, r); err != nil {
			return err
		}

		return nil
	}()

	return r, svc.recordAction(ctx, raProps, ProjectTaskActionUpdate, err)
}

func (svc projectTask) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		raProps = &projectTaskActionProps{}
		r       *types.ProjectTask
	)

	err = func() (err error) {
		if r, err = loadProjectTask(ctx, svc.store, ID); err != nil {
			return
		}

		raProps.setProjectTask(r)

		if !svc.ac.CanDeleteProjectTask(ctx, r) {
			return ProjectTaskErrNotAllowedToDelete()
		}

		r.DeletedAt = now()
		if err = store.UpdateProjectTask(ctx, svc.store, r); err != nil {
			return
		}

		return nil
	}()

	return svc.recordAction(ctx, raProps, ProjectTaskActionDelete, err)
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
