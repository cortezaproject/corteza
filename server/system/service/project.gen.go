package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"

	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/pkg/scope"
	"github.com/crusttech/human/server/store"
	types "github.com/crusttech/human/server/system/types"
)

type project struct {
	actionlog actionlog.Recorder
	store     store.Storer
	ac        projectAccessController
}

func Project() *project {
	return &project{
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		ac:        DefaultAccessControl,
	}
}

type projectServices struct {
	scope scope.Scope
	caps  scope.Capabilities
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

func (svc *project) checkScope(ctx context.Context, cap scope.Capability) error {
	if err := scope.RequireTenantMembership(ctx); err != nil {
		return err
	}
	return scope.RequireCapability(ctx, cap)
}

func (svc *project) scopeServices(ctx context.Context) *projectServices {
	return &projectServices{
		scope: scope.GetScopeFromContext(ctx),
		caps:  scope.GetCapabilitiesFromContext(ctx),
	}
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
