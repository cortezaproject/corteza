package service

import (
	"context"
	"strconv"
	"time"

	composeTypes "github.com/crusttech/human/server/compose/types"
	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/handle"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

type (
	// project struct + the Project() constructor are generated into project.gen.go.

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
			new.Status = types.ProjectStatusDraft
		}
		if err = validateProjectStatus(new.Status); err != nil {
			return err
		}

		if new.Config.Mode == "" {
			new.Config.Mode = types.ProjectModeFree
		}
		if !new.Config.Mode.Valid() {
			return ProjectErrInvalidMode()
		}

		// FRIA requirement is derived from the deployer answers once, at
		// creation, so the pipeline shape doesn't silently change later.
		new.Config.FriaRequired = friaRequired(new.Config.DeployerCategories)

		if err = svc.uniqueCheck(ctx, new); err != nil {
			return err
		}

		new.ID = nextID()
		new.CreatedAt = *now()
		new.CreatedBy = a.GetIdentityFromContext(ctx).Identity()

		// Project-scoped compose resources (modules, pages, charts) live in a
		// namespace created alongside the project. Projects are referenced by
		// ID; handle and slug stay empty unless a client explicitly sets one.
		// The scope refs are stamped explicitly: the request context carries
		// no project scope yet at creation time.
		ns := &composeTypes.Namespace{
			ID:        nextID(),
			TenantID:  new.TenantID,
			ProjectID: new.ID,
			Slug:      new.Handle,
			Name:      new.Meta.Short,
			Enabled:   true,
			CreatedAt: *now(),
		}
		if ns.Name == "" {
			ns.Name = "project-" + strconv.FormatUint(new.ID, 10)
		}
		new.Config.NamespaceID = ns.ID

		err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
			// Slug uniqueness only matters when one is actually set; empty
			// slugs are allowed to repeat (partial unique index).
			if ns.Slug != "" {
				if existing, e := store.LookupComposeNamespaceBySlug(ctx, s, ns.Slug); e == nil && existing != nil {
					return ProjectErrHandleNotUnique()
				} else if e != nil && !errors.IsNotFound(e) {
					return e
				}
			}

			if err = store.CreateComposeNamespace(ctx, s, ns); err != nil {
				return err
			}

			return store.CreateProject(ctx, s, new)
		})
		if err != nil {
			return
		}

		if err = label.Create(ctx, svc.store, new); err != nil {
			return
		}

		// The creator becomes the first member: in gated mode they drive the
		// pipeline as a developer.
		creator := &types.ProjectMember{
			ID:         nextID(),
			ProjectID:  new.ID,
			TenantID:   new.TenantID,
			UserID:     new.CreatedBy,
			RolePreset: types.ProjectRoleDeveloper,
			InvitedBy:  new.CreatedBy,
			CreatedAt:  *now(),
		}
		if creator.UserID != 0 {
			if err = store.CreateProjectMember(ctx, svc.store, creator); err != nil {
				return
			}
		}

		p = new
		return nil
	}()

	return p, svc.recordAction(ctx, paProps, ProjectActionCreate, err)
}

func (svc *project) Update(ctx context.Context, upd *types.Project) (p *types.Project, err error) {
	var (
		paProps  = &projectActionProps{update: upd}
		existing *types.Project
	)

	err = func() (err error) {
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

		// Mode and namespace binding are immutable; FRIA derivation is fixed at
		// creation. Governance state only changes through the dedicated
		// governance operations.
		upd.Config.Mode = existing.Config.Mode
		upd.Config.NamespaceID = existing.Config.NamespaceID
		upd.Config.DeployerCategories = existing.Config.DeployerCategories
		upd.Config.FriaRequired = existing.Config.FriaRequired
		upd.Governance = existing.Governance

		if err = store.UpdateProject(ctx, svc.store, upd); err != nil {
			return
		}

		if err = label.Update(ctx, svc.store, upd); err != nil {
			return
		}

		p = upd
		return nil
	}()

	return p, svc.recordAction(ctx, paProps, ProjectActionUpdate, err, existing, upd)
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

		// The project's namespace follows its lifecycle — leaving it active
		// would squat the slug and block future creations.
		return store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
			if err := store.UpdateProject(ctx, s, p); err != nil {
				return err
			}
			return svc.setNamespaceDeleted(ctx, s, p.Config.NamespaceID, p.DeletedAt)
		})
	}()

	return svc.recordAction(ctx, paProps, ProjectActionDelete, err)
}

// setNamespaceDeleted stamps (or clears) DeletedAt on the project's compose
// namespace. A missing namespace is not an error — older projects may predate
// the auto-created namespace.
func (svc *project) setNamespaceDeleted(ctx context.Context, s store.Storer, namespaceID uint64, deletedAt *time.Time) error {
	if namespaceID == 0 {
		return nil
	}

	ns, err := store.LookupComposeNamespaceByID(ctx, s, namespaceID)
	if errors.IsNotFound(err) {
		return nil
	} else if err != nil {
		return err
	}

	ns.DeletedAt = deletedAt
	return store.UpdateComposeNamespace(ctx, s, ns)
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

		return store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
			if err := store.UpdateProject(ctx, s, p); err != nil {
				return err
			}
			return svc.setNamespaceDeleted(ctx, s, p.Config.NamespaceID, nil)
		})
	}()

	return svc.recordAction(ctx, paProps, ProjectActionUndelete, err)
}

// --- members ---

func (svc *project) onSearchMembers(ctx context.Context, _ *projectActionProps, filter types.ProjectMemberFilter) (set types.ProjectMemberSet, f types.ProjectMemberFilter, err error) {
	var p *types.Project
	if p, err = loadProject(ctx, svc.store, filter.ProjectID); err != nil {
		return
	}

	if !svc.ac.CanReadProject(ctx, p) {
		return set, f, ProjectErrNotAllowedToRead()
	}

	return store.SearchProjectMembers(ctx, svc.store, filter)
}

func (svc *project) onAddMember(ctx context.Context, _ *projectActionProps, m *types.ProjectMember) (res *types.ProjectMember, err error) {
	var p *types.Project
	if p, err = loadProject(ctx, svc.store, m.ProjectID); err != nil {
		return
	}

	if !svc.ac.CanManageMembersOnProject(ctx, p) {
		return nil, ProjectErrNotAllowedToManageMembers()
	}

	if m.UserID == 0 {
		return nil, ProjectErrMemberNotFound()
	}

	if m.RolePreset == "" {
		m.RolePreset = p.Config.DefaultMemberRole
	}
	if m.RolePreset == "" {
		m.RolePreset = types.ProjectRoleMember
	}
	if !m.RolePreset.Valid() {
		return nil, ProjectErrInvalidRolePreset()
	}

	// one membership record per user per project
	if existing, e := store.LookupProjectMemberByProjectIDUserID(ctx, svc.store, m.ProjectID, m.UserID); e == nil && existing != nil {
		return nil, ProjectErrMemberAlreadyExists()
	} else if e != nil && !errors.IsNotFound(e) {
		return nil, e
	}

	m.ID = nextID()
	m.TenantID = p.TenantID
	m.CreatedAt = *now()
	m.InvitedBy = a.GetIdentityFromContext(ctx).Identity()

	if err = store.CreateProjectMember(ctx, svc.store, m); err != nil {
		return
	}

	return m, nil
}

func (svc *project) onUpdateMember(ctx context.Context, _ *projectActionProps, m *types.ProjectMember) (res *types.ProjectMember, err error) {
	var p *types.Project
	if p, err = loadProject(ctx, svc.store, m.ProjectID); err != nil {
		return
	}

	if !svc.ac.CanManageMembersOnProject(ctx, p) {
		return nil, ProjectErrNotAllowedToManageMembers()
	}

	var existing *types.ProjectMember
	if existing, err = store.LookupProjectMemberByProjectIDUserID(ctx, svc.store, m.ProjectID, m.UserID); errors.IsNotFound(err) {
		return nil, ProjectErrMemberNotFound()
	} else if err != nil {
		return nil, err
	}

	if !m.RolePreset.Valid() {
		return nil, ProjectErrInvalidRolePreset()
	}

	existing.RolePreset = m.RolePreset
	existing.UpdatedAt = now()

	if err = store.UpdateProjectMember(ctx, svc.store, existing); err != nil {
		return
	}

	return existing, nil
}

func (svc *project) onRemoveMember(ctx context.Context, _ *projectActionProps, projectID, userID uint64) (err error) {
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
	case types.ProjectStatusDraft,
		types.ProjectStatusActive,
		types.ProjectStatusPublished,
		types.ProjectStatusArchived,
		types.ProjectStatusSuspended:
		return nil
	}
	return ProjectErrInvalidStatus()
}

// friaRequired: any positive deployer-category answer makes a Fundamental
// Rights Impact Assessment step mandatory (EU AI Act Art. 27).
func friaRequired(d types.ProjectDeployerCategories) bool {
	return d.PublicAuthorityAnnex3 || d.PrivateEssentialServices || d.InsuranceBanking
}

// toLabeledProjects is generated into project.gen.go.

func loadProject(ctx context.Context, s store.Projects, ID uint64) (res *types.Project, err error) {
	if ID == 0 {
		return nil, ProjectErrInvalidID()
	}

	if res, err = store.LookupProjectByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, ProjectErrNotFound()
	}

	return
}
