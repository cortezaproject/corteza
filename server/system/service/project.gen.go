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
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/pkg/scope"
	"github.com/crusttech/human/server/store"
	types "github.com/crusttech/human/server/system/types"
)

type project struct {
	actionlog actionlog.Recorder
	store     store.Storer
	ac        projectAccessController
	services  *projectServices
}

func (svc *project) FindByID(ctx context.Context, ID uint64) (res *types.Project, err error) {
	var (
		aProps = &projectActionProps{project: &types.Project{ID: ID}}
	)

	err = func() error {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}
		if res, err = loadProject(ctx, svc.store, ID); err != nil {
			return ProjectErrInvalidID().Wrap(err)
		}

		if res, err = svc.afterLookup(ctx, res); err != nil {
			return err
		}

		aProps.setProject(res)

		if err = svc.guard(ctx, res); err != nil {
			return err
		}

		if !svc.ac.CanReadProject(ctx, res) {
			return ProjectErrNotAllowedToRead()
		}

		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ProjectActionLookup, err)
}

func (svc *project) Search(ctx context.Context, filter types.ProjectFilter) (set types.ProjectSet, f types.ProjectFilter, err error) {
	var (
		aProps = &projectActionProps{filter: &filter}
	)

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *types.Project) (bool, error) {
		if !svc.ac.CanReadProject(ctx, res) {
			return false, nil
		}

		return true, nil
	}

	err = func() error {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}
		if !svc.ac.CanSearchProjects(ctx) {
			return ProjectErrNotAllowedToSearch()
		}

		if len(filter.Labels) > 0 {
			filter.LabeledIDs, err = label.Search(
				ctx,
				svc.store,
				types.Project{}.LabelResourceKind(),
				filter.Labels,
			)
			if err != nil {
				return err
			}

			// labels specified but no labeled resources found
			if len(filter.LabeledIDs) == 0 {
				return nil
			}
		}

		if set, f, err = store.SearchProjects(ctx, svc.store, filter); err != nil {
			return err
		}

		if err = label.Load(ctx, svc.store, toLabeledProjects(set)...); err != nil {
			return err
		}

		return nil
	}()

	return set, f, svc.recordAction(ctx, aProps, ProjectActionSearch, err)
}

func (svc *project) Create(ctx context.Context, new *types.Project) (res *types.Project, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &projectActionProps{project: new, new: new}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}
		if !svc.ac.CanCreateProject(ctx) {
			return ProjectErrNotAllowedToCreate()
		}
		res = new
		return svc.onCreate(ctx, new)
	}()

	return res, svc.recordAction(ctx, aProps, ProjectActionCreate, err)
}

func (svc *project) Update(ctx context.Context, upd *types.Project) (res *types.Project, err error) {
	var (
		aProps = &projectActionProps{update: upd}
		old    *types.Project
	)
	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}
		if res, err = loadProject(ctx, svc.store, upd.ID); err != nil {
			return
		}

		old = res.Clone()
		aProps.setProject(res)

		if err = svc.guard(ctx, res); err != nil {
			return err
		}

		if upd.Handle != res.Handle && !handle.IsValid(upd.Handle) {
			return ProjectErrInvalidHandle()
		}

		if !svc.ac.CanUpdateProject(ctx, res) {
			return ProjectErrNotAllowedToUpdate()
		}

		// Test if stale (update has an older version of data)
		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return ProjectErrStaleData()
		}

		if err = svc.beforeUpdate(ctx, upd, res); err != nil {
			return err
		}
		res.Handle = upd.Handle
		res.Config = upd.Config
		res.Meta = upd.Meta
		res.UpdatedBy = upd.UpdatedBy
		res.UpdatedAt = now()

		if err = store.UpdateProject(ctx, svc.store, res); err != nil {
			return err
		}

		if label.Changed(res.Labels, upd.Labels) {
			if err = label.Update(ctx, svc.store, upd); err != nil {
				return
			}
			res.Labels = upd.Labels
		}
		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ProjectActionUpdate, err, old, res)
}

func (svc *project) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &projectActionProps{}
		res    *types.Project
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadProject(ctx, s, ID); err != nil {
			return
		}

		if err = label.Load(ctx, svc.store, res); err != nil {
			return err
		}

		aProps.setProject(res)

		if err = svc.guard(ctx, res); err != nil {
			return
		}

		return svc.onDelete(ctx, s, res, aProps)
	})

	return svc.recordAction(ctx, aProps, ProjectActionDelete, err)
}

func (svc *project) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &projectActionProps{}
		res    *types.Project
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadProject(ctx, s, ID); err != nil {
			return
		}

		if err = label.Load(ctx, svc.store, res); err != nil {
			return err
		}

		aProps.setProject(res)

		if err = svc.guard(ctx, res); err != nil {
			return
		}

		return svc.onUndelete(ctx, s, res, aProps)
	})

	return svc.recordAction(ctx, aProps, ProjectActionUndelete, err)
}

func loadProject(ctx context.Context, s store.Projects, ID uint64) (res *types.Project, err error) {
	if ID == 0 {
		return nil, ProjectErrInvalidID()
	}

	if res, err = store.LookupProjectByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, ProjectErrNotFound()
	}

	return
}

// toLabeledProjects converts to []label.LabeledResource
func toLabeledProjects(set []*types.Project) []label.LabeledResource {
	if len(set) == 0 {
		return nil
	}

	ll := make([]label.LabeledResource, len(set))
	for i := range set {
		ll[i] = set[i]
	}

	return ll
}
func (svc *project) guard(_ context.Context, _ *types.Project) error { return nil }

func (svc *project) checkScope(ctx context.Context, cap scope.Capability) error {
	if err := scope.RequireTenantMembership(ctx); err != nil {
		return err
	}
	return scope.RequireCapability(ctx, cap)
}

func (svc *project) SearchMembers(ctx context.Context, filter types.ProjectMemberFilter) (set types.ProjectMemberSet, f types.ProjectMemberFilter, err error) {
	var (
		aProps = &projectActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}

		set, f, err = svc.onSearchMembers(ctx, aProps, filter)
		return err
	}()

	return set, f, svc.recordAction(ctx, aProps, ProjectActionSearchMembers, err)
}

func (svc *project) AddMember(ctx context.Context, m *types.ProjectMember) (res *types.ProjectMember, err error) {
	var (
		aProps = &projectActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		res, err = svc.onAddMember(ctx, aProps, m)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, ProjectActionAddMember, err)
}

func (svc *project) UpdateMember(ctx context.Context, m *types.ProjectMember) (res *types.ProjectMember, err error) {
	var (
		aProps = &projectActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		res, err = svc.onUpdateMember(ctx, aProps, m)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, ProjectActionUpdateMember, err)
}

func (svc *project) RemoveMember(ctx context.Context, projectID uint64, userID uint64) (err error) {
	var (
		aProps = &projectActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		err = svc.onRemoveMember(ctx, aProps, projectID, userID)
		return err
	}()

	return svc.recordAction(ctx, aProps, ProjectActionRemoveMember, err)
}
