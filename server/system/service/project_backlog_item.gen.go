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

type projectBacklogItemAccessController interface {
	CanCreateProjectBacklogItem(context.Context) bool
	CanSearchProjectBacklogItems(context.Context) bool
	CanReadProjectBacklogItem(context.Context, *types.ProjectBacklogItem) bool
	CanUpdateProjectBacklogItem(context.Context, *types.ProjectBacklogItem) bool
	CanDeleteProjectBacklogItem(context.Context, *types.ProjectBacklogItem) bool
}

type projectBacklogItem struct {
	actionlog actionlog.Recorder
	store     store.Storer
	ac        projectBacklogItemAccessController
}

type projectBacklogItemServices struct {
	scope scope.Scope
	caps  scope.Capabilities
}

func (svc *projectBacklogItem) FindByID(ctx context.Context, ID uint64) (res *types.ProjectBacklogItem, err error) {
	var (
		aProps = &projectBacklogItemActionProps{projectBacklogItem: &types.ProjectBacklogItem{ID: ID}}
	)

	err = func() error {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}
		if res, err = loadProjectBacklogItem(ctx, svc.store, ID); err != nil {
			return ProjectBacklogItemErrInvalidID().Wrap(err)
		}

		aProps.setProjectBacklogItem(res)

		if err = svc.guard(ctx, res); err != nil {
			return err
		}

		if !svc.ac.CanReadProjectBacklogItem(ctx, res) {
			return ProjectBacklogItemErrNotAllowedToRead()
		}

		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ProjectBacklogItemActionLookup, err)
}

func (svc *projectBacklogItem) Search(ctx context.Context, filter types.ProjectBacklogItemFilter) (set types.ProjectBacklogItemSet, f types.ProjectBacklogItemFilter, err error) {
	var (
		aProps = &projectBacklogItemActionProps{filter: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.ProjectBacklogItem) (bool, error) {
		if !svc.ac.CanReadProjectBacklogItem(ctx, res) {
			return false, nil
		}

		return true, nil
	}

	err = func() error {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}
		if !svc.ac.CanSearchProjectBacklogItems(ctx) {
			return ProjectBacklogItemErrNotAllowedToSearch()
		}

		if set, f, err = store.SearchProjectBacklogItems(ctx, svc.store, filter); err != nil {
			return err
		}

		return nil
	}()

	return set, f, svc.recordAction(ctx, aProps, ProjectBacklogItemActionSearch, err)
}

func (svc *projectBacklogItem) Create(ctx context.Context, new *types.ProjectBacklogItem) (res *types.ProjectBacklogItem, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &projectBacklogItemActionProps{projectBacklogItem: new, new: new}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}
		if !svc.ac.CanCreateProjectBacklogItem(ctx) {
			return ProjectBacklogItemErrNotAllowedToCreate()
		}

		if err = svc.beforeCreate(ctx, new); err != nil {
			return err
		}

		new.ID = nextID()
		new.CreatedAt = *now()

		if err = store.CreateProjectBacklogItem(ctx, svc.store, new); err != nil {
			return
		}

		res = new
		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ProjectBacklogItemActionCreate, err)
}

func (svc *projectBacklogItem) Update(ctx context.Context, upd *types.ProjectBacklogItem) (res *types.ProjectBacklogItem, err error) {
	var (
		aProps = &projectBacklogItemActionProps{update: upd}
		old    *types.ProjectBacklogItem
	)
	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}
		if res, err = loadProjectBacklogItem(ctx, svc.store, upd.ID); err != nil {
			return
		}

		old = res.Clone()
		aProps.setProjectBacklogItem(res)

		if err = svc.guard(ctx, res); err != nil {
			return err
		}

		if !svc.ac.CanUpdateProjectBacklogItem(ctx, res) {
			return ProjectBacklogItemErrNotAllowedToUpdate()
		}

		// Test if stale (update has an older version of data)
		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return ProjectBacklogItemErrStaleData()
		}

		if err = svc.beforeUpdate(ctx, upd, res); err != nil {
			return err
		}
		res.RevisionID = upd.RevisionID
		res.Title = upd.Title
		res.Description = upd.Description
		res.Category = upd.Category
		res.EventID = upd.EventID
		res.Assignee = upd.Assignee
		res.Priority = upd.Priority
		res.Status = upd.Status
		res.DateDue = upd.DateDue
		res.UpdatedAt = now()

		if err = store.UpdateProjectBacklogItem(ctx, svc.store, res); err != nil {
			return err
		}
		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ProjectBacklogItemActionUpdate, err, old, res)
}

func (svc *projectBacklogItem) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &projectBacklogItemActionProps{}
		res    *types.ProjectBacklogItem
	)
	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		if res, err = loadProjectBacklogItem(ctx, svc.store, ID); err != nil {
			return
		}

		if err = svc.guard(ctx, res); err != nil {
			return err
		}

		aProps.setProjectBacklogItem(res)

		if !svc.ac.CanDeleteProjectBacklogItem(ctx, res) {
			return ProjectBacklogItemErrNotAllowedToDelete()
		}

		if err = svc.beforeDelete(ctx, res); err != nil {
			return err
		}

		res.DeletedAt = now()
		if err = store.UpdateProjectBacklogItem(ctx, svc.store, res); err != nil {
			return
		}
		return nil
	}()

	return svc.recordAction(ctx, aProps, ProjectBacklogItemActionDelete, err)
}

func loadProjectBacklogItem(ctx context.Context, s store.ProjectBacklogItems, ID uint64) (res *types.ProjectBacklogItem, err error) {
	if ID == 0 {
		return nil, ProjectBacklogItemErrInvalidID()
	}

	if res, err = store.LookupProjectBacklogItemByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, ProjectBacklogItemErrNotFound()
	}

	return
}
func (svc *projectBacklogItem) guard(_ context.Context, _ *types.ProjectBacklogItem) error {
	return nil
}

func (svc *projectBacklogItem) checkScope(ctx context.Context, cap scope.Capability) error {
	if err := scope.RequireTenantMembership(ctx); err != nil {
		return err
	}
	return scope.RequireCapability(ctx, cap)
}
func (svc *projectBacklogItem) scopeServices(ctx context.Context) *projectBacklogItemServices {
	return &projectBacklogItemServices{
		scope: scope.GetScopeFromContext(ctx),
		caps:  scope.GetCapabilitiesFromContext(ctx),
	}
}
