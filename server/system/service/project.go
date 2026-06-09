package service

import (
	"context"

	"github.com/crusttech/human/server/pkg/actionlog"
	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/handle"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

type (
	project struct {
		actionlog actionlog.Recorder
		store     store.Storer
		ac        projectAccessController
	}

	projectAccessController interface {
		CanCreateProject(ctx context.Context) bool
		CanSearchProjects(ctx context.Context) bool
		CanReadProject(ctx context.Context, p *types.Project) bool
		CanUpdateProject(ctx context.Context, p *types.Project) bool
		CanDeleteProject(ctx context.Context, p *types.Project) bool
		CanManageMembersOnProject(ctx context.Context, p *types.Project) bool
	}

	ProjectService interface {
		FindByID(ctx context.Context, ID uint64) (*types.Project, error)
		FindByHandle(ctx context.Context, handle string) (*types.Project, error)
		Search(ctx context.Context, filter types.ProjectFilter) (types.ProjectSet, types.ProjectFilter, error)

		Create(ctx context.Context, new *types.Project) (*types.Project, error)
		Update(ctx context.Context, upd *types.Project) (*types.Project, error)
		DeleteByID(ctx context.Context, ID uint64) error
		UndeleteByID(ctx context.Context, ID uint64) error

		SearchMembers(ctx context.Context, filter types.ProjectMemberFilter) (types.ProjectMemberSet, types.ProjectMemberFilter, error)
		AddMember(ctx context.Context, m *types.ProjectMember) (*types.ProjectMember, error)
		UpdateMember(ctx context.Context, m *types.ProjectMember) (*types.ProjectMember, error)
		RemoveMember(ctx context.Context, projectID, userID uint64) error
	}
)

func Project() *project {
	return &project{
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
	}
}

func (svc *project) FindByID(ctx context.Context, ID uint64) (p *types.Project, err error) {
	var paProps = &projectActionProps{project: &types.Project{ID: ID}}

	err = func() error {
		if p, err = loadProject(ctx, svc.store, ID); err != nil {
			return err
		}

		paProps.setProject(p)

		if !svc.ac.CanReadProject(ctx, p) {
			return ProjectErrNotAllowedToRead()
		}

		if err = label.Load(ctx, svc.store, p); err != nil {
			return err
		}

		return nil
	}()

	return p, svc.recordAction(ctx, paProps, ProjectActionLookup, err)
}

func (svc *project) FindByHandle(ctx context.Context, h string) (p *types.Project, err error) {
	var paProps = &projectActionProps{project: &types.Project{Handle: h}}

	err = func() error {
		if p, err = store.LookupProjectByHandle(ctx, svc.store, h); errors.IsNotFound(err) {
			return ProjectErrNotFound()
		} else if err != nil {
			return err
		}

		paProps.setProject(p)

		if !svc.ac.CanReadProject(ctx, p) {
			return ProjectErrNotAllowedToRead()
		}

		if err = label.Load(ctx, svc.store, p); err != nil {
			return err
		}

		return nil
	}()

	return p, svc.recordAction(ctx, paProps, ProjectActionLookup, err)
}

func (svc *project) Create(ctx context.Context, new *types.Project) (p *types.Project, err error) {
	var paProps = &projectActionProps{new: new}

	err = func() (err error) {
		if !svc.ac.CanCreateProject(ctx) {
			return ProjectErrNotAllowedToCreate()
		}

		if !handle.IsValid(new.Handle) {
			return ProjectErrInvalidHandle()
		}

		if new.Status == "" {
			new.Status = types.ProjectStatusActive
		}
		if err = validateProjectStatus(new.Status); err != nil {
			return err
		}

		if err = svc.uniqueCheck(ctx, new); err != nil {
			return err
		}

		new.ID = nextID()
		new.CreatedAt = *now()
		new.CreatedBy = a.GetIdentityFromContext(ctx).Identity()

		if err = store.CreateProject(ctx, svc.store, new); err != nil {
			return
		}

		if err = label.Create(ctx, svc.store, new); err != nil {
			return
		}

		p = new
		return nil
	}()

	return p, svc.recordAction(ctx, paProps, ProjectActionCreate, err)
}

func (svc *project) Update(ctx context.Context, upd *types.Project) (p *types.Project, err error) {
	var paProps = &projectActionProps{update: upd}

	err = func() (err error) {
		var existing *types.Project
		if existing, err = loadProject(ctx, svc.store, upd.ID); err != nil {
			return
		}

		paProps.setProject(existing)

		if !svc.ac.CanUpdateProject(ctx, existing) {
			return ProjectErrNotAllowedToUpdate()
		}

		if isStale(upd.UpdatedAt, existing.UpdatedAt, existing.CreatedAt) {
			return ProjectErrStaleData()
		}

		if !handle.IsValid(upd.Handle) {
			return ProjectErrInvalidHandle()
		}

		if upd.Status == "" {
			upd.Status = existing.Status
		}
		if err = validateProjectStatus(upd.Status); err != nil {
			return err
		}

		if upd.Handle != existing.Handle {
			if err = svc.uniqueCheck(ctx, upd); err != nil {
				return err
			}
		}

		upd.UpdatedAt = now()
		upd.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()
		upd.CreatedAt = existing.CreatedAt
		upd.CreatedBy = existing.CreatedBy
		upd.DeletedAt = existing.DeletedAt
		upd.TenantID = existing.TenantID

		if err = store.UpdateProject(ctx, svc.store, upd); err != nil {
			return
		}

		if err = label.Update(ctx, svc.store, upd); err != nil {
			return
		}

		p = upd
		return nil
	}()

	return p, svc.recordAction(ctx, paProps, ProjectActionUpdate, err)
}

func (svc *project) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var paProps = &projectActionProps{project: &types.Project{ID: ID}}

	err = func() (err error) {
		var p *types.Project
		if p, err = loadProject(ctx, svc.store, ID); err != nil {
			return
		}

		paProps.setProject(p)

		if !svc.ac.CanDeleteProject(ctx, p) {
			return ProjectErrNotAllowedToDelete()
		}

		p.DeletedAt = now()
		p.DeletedBy = a.GetIdentityFromContext(ctx).Identity()
		if err = store.UpdateProject(ctx, svc.store, p); err != nil {
			return
		}

		return nil
	}()

	return svc.recordAction(ctx, paProps, ProjectActionDelete, err)
}

func (svc *project) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	var paProps = &projectActionProps{project: &types.Project{ID: ID}}

	err = func() (err error) {
		var p *types.Project
		if p, err = loadProject(ctx, svc.store, ID); err != nil {
			return
		}

		paProps.setProject(p)

		if !svc.ac.CanDeleteProject(ctx, p) {
			return ProjectErrNotAllowedToDelete()
		}

		p.DeletedAt = nil
		p.DeletedBy = 0
		if err = store.UpdateProject(ctx, svc.store, p); err != nil {
			return
		}

		return nil
	}()

	return svc.recordAction(ctx, paProps, ProjectActionUndelete, err)
}

func (svc *project) Search(ctx context.Context, filter types.ProjectFilter) (set types.ProjectSet, f types.ProjectFilter, err error) {
	var paProps = &projectActionProps{search: &filter}

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

	return set, f, svc.recordAction(ctx, paProps, ProjectActionSearch, err)
}

// --- members ---

func (svc *project) SearchMembers(ctx context.Context, filter types.ProjectMemberFilter) (set types.ProjectMemberSet, f types.ProjectMemberFilter, err error) {
	err = func() error {
		var p *types.Project
		if p, err = loadProject(ctx, svc.store, filter.ProjectID); err != nil {
			return err
		}

		if !svc.ac.CanReadProject(ctx, p) {
			return ProjectErrNotAllowedToRead()
		}

		if set, f, err = store.SearchProjectMembers(ctx, svc.store, filter); err != nil {
			return err
		}

		return nil
	}()

	return set, f, err
}

func (svc *project) AddMember(ctx context.Context, m *types.ProjectMember) (res *types.ProjectMember, err error) {
	err = func() (err error) {
		var p *types.Project
		if p, err = loadProject(ctx, svc.store, m.ProjectID); err != nil {
			return
		}

		if !svc.ac.CanManageMembersOnProject(ctx, p) {
			return ProjectErrNotAllowedToManageMembers()
		}

		if m.UserID == 0 {
			return ProjectErrMemberNotFound()
		}

		if m.RolePreset == "" {
			m.RolePreset = p.Config.DefaultMemberRole
		}
		if m.RolePreset == "" {
			m.RolePreset = types.ProjectRoleMember
		}
		if !m.RolePreset.Valid() {
			return ProjectErrInvalidRolePreset()
		}

		// one membership record per user per project
		if existing, e := store.LookupProjectMemberByProjectIDUserID(ctx, svc.store, m.ProjectID, m.UserID); e == nil && existing != nil {
			return ProjectErrMemberAlreadyExists()
		} else if e != nil && !errors.IsNotFound(e) {
			return e
		}

		m.ID = nextID()
		m.TenantID = p.TenantID
		m.CreatedAt = *now()
		m.InvitedBy = a.GetIdentityFromContext(ctx).Identity()

		if err = store.CreateProjectMember(ctx, svc.store, m); err != nil {
			return
		}

		res = m
		return nil
	}()

	return res, err
}

func (svc *project) UpdateMember(ctx context.Context, m *types.ProjectMember) (res *types.ProjectMember, err error) {
	err = func() (err error) {
		var p *types.Project
		if p, err = loadProject(ctx, svc.store, m.ProjectID); err != nil {
			return
		}

		if !svc.ac.CanManageMembersOnProject(ctx, p) {
			return ProjectErrNotAllowedToManageMembers()
		}

		var existing *types.ProjectMember
		if existing, err = store.LookupProjectMemberByProjectIDUserID(ctx, svc.store, m.ProjectID, m.UserID); errors.IsNotFound(err) {
			return ProjectErrMemberNotFound()
		} else if err != nil {
			return err
		}

		if !m.RolePreset.Valid() {
			return ProjectErrInvalidRolePreset()
		}

		existing.RolePreset = m.RolePreset
		existing.UpdatedAt = now()

		if err = store.UpdateProjectMember(ctx, svc.store, existing); err != nil {
			return
		}

		res = existing
		return nil
	}()

	return res, err
}

func (svc *project) RemoveMember(ctx context.Context, projectID, userID uint64) (err error) {
	err = func() (err error) {
		var p *types.Project
		if p, err = loadProject(ctx, svc.store, projectID); err != nil {
			return
		}

		if !svc.ac.CanManageMembersOnProject(ctx, p) {
			return ProjectErrNotAllowedToManageMembers()
		}

		var existing *types.ProjectMember
		if existing, err = store.LookupProjectMemberByProjectIDUserID(ctx, svc.store, projectID, userID); errors.IsNotFound(err) {
			return ProjectErrMemberNotFound()
		} else if err != nil {
			return err
		}

		existing.DeletedAt = now()
		return store.UpdateProjectMember(ctx, svc.store, existing)
	}()

	return err
}

// --- helpers ---

func (svc *project) uniqueCheck(ctx context.Context, p *types.Project) error {
	if p.Handle == "" {
		return nil
	}
	if e, err := store.LookupProjectByHandle(ctx, svc.store, p.Handle); err == nil && e != nil && e.ID != p.ID {
		return ProjectErrHandleNotUnique()
	} else if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

func validateProjectStatus(s types.ProjectStatus) error {
	switch s {
	case types.ProjectStatusActive, types.ProjectStatusArchived, types.ProjectStatusSuspended:
		return nil
	}
	return ProjectErrInvalidStatus()
}

func toLabeledProjects(set types.ProjectSet) []label.LabeledResource {
	if len(set) == 0 {
		return nil
	}
	ll := make([]label.LabeledResource, len(set))
	for i := range set {
		ll[i] = set[i]
	}
	return ll
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
