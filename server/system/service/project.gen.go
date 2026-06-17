package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"

	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/store"
	types "github.com/crusttech/human/server/system/types"
)

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
		if !svc.ac.CanCreateProject(ctx) {
			return ProjectErrNotAllowedToCreate()
		}

		if err = svc.beforeCreate(ctx, new); err != nil {
			return err
		}

		new.ID = nextID()
		new.CreatedAt = *now()

		if err = store.CreateProject(ctx, svc.store, new); err != nil {
			return
		}

		if err = label.Create(ctx, svc.store, new); err != nil {
			return
		}

		res = new
		return nil
	}()

	return res, svc.recordAction(ctx, aProps, ProjectActionCreate, err)
}

func (svc *project) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &projectActionProps{}
		res    *types.Project
	)

	err = func() (err error) {
		if res, err = loadProject(ctx, svc.store, ID); err != nil {
			return
		}

		aProps.setProject(res)

		if !svc.ac.CanDeleteProject(ctx, res) {
			return ProjectErrNotAllowedToDelete()
		}

		if err = svc.beforeDelete(ctx, res); err != nil {
			return err
		}

		res.DeletedAt = now()
		if err = store.UpdateProject(ctx, svc.store, res); err != nil {
			return
		}
		return nil
	}()

	return svc.recordAction(ctx, aProps, ProjectActionDelete, err)
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
