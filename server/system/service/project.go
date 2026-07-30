package service

import (
	"context"
	"strconv"
	"time"

	composeTypes "github.com/crusttech/human/server/compose/types"
	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/dal"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/handle"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/pkg/scope"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

type (
	projectNamespaceSvc interface {
		CloneFromStore(ctx context.Context, sourceNsID uint64, dup *composeTypes.Namespace) (*composeTypes.Namespace, error)
	}

	projectDALImporter interface {
		RunImportForeground(ctx context.Context, mappingID uint64) error
	}

	projectDALConnSvc interface {
		ReplaceConnection(ctx context.Context, conn *dal.ConnectionWrap, isDefault bool) error
		RemoveConnection(ctx context.Context, ID uint64) error
	}

	projectRecordSvc interface {
		Bulk(ctx context.Context, skipFailed bool, oo ...*composeTypes.RecordBulkOperation) ([]composeTypes.RecordBulkOperationResult, error)
		Search(ctx context.Context, f composeTypes.RecordFilter) (composeTypes.RecordSet, composeTypes.RecordFilter, error)
	}
)

func Project() *project {
	return &project{
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		ac:        DefaultAccessControl,
		services:  &projectServices{},
	}
}

type (
	projectAccessController interface {
		CanCreateProject(ctx context.Context) bool
		CanSearchProjects(ctx context.Context) bool
		CanReadProject(ctx context.Context, p *types.Project) bool
		CanUpdateProject(ctx context.Context, p *types.Project) bool
		CanDeleteProject(ctx context.Context, p *types.Project) bool
		CanManageMembersOnProject(ctx context.Context, p *types.Project) bool
		CanReviseProject(ctx context.Context, p *types.Project) bool
		CanPublishProject(ctx context.Context, p *types.Project) bool
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

		CreateRevision(ctx context.Context, projectID uint64) (*types.Project, error)
		ListRevisions(ctx context.Context, projectID uint64) (types.ProjectSet, error)
		DeploymentPlan(ctx context.Context, projectID uint64) (*types.ProjectDeploymentPlan, error)
		Publish(ctx context.Context, projectID uint64, req types.PublishRequest) (*types.Project, error)
	}
)

// projectServices extends the generated scope/caps struct with project-specific
// dependencies wired at construction time.
type projectServices struct {
	scope     scope.Scope
	caps      scope.Capabilities
	nsSvc     projectNamespaceSvc
	dalSvc    projectDALImporter
	dalConns  projectDALConnSvc
	recordSvc projectRecordSvc
}

func (svc *project) scopeServices(ctx context.Context) *projectServices {
	return &projectServices{
		scope:     scope.GetScopeFromContext(ctx),
		caps:      scope.GetCapabilitiesFromContext(ctx),
		nsSvc:     svc.services.nsSvc,
		dalSvc:    svc.services.dalSvc,
		dalConns:  svc.services.dalConns,
		recordSvc: svc.services.recordSvc,
	}
}

func (svc *project) afterLookup(ctx context.Context, p *types.Project) (*types.Project, error) {
	if err := label.Load(ctx, svc.store, p); err != nil {
		return nil, err
	}
	return p, nil
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

func (svc *project) onCreate(ctx context.Context, new *types.Project) error {
	if !handle.IsValid(new.Handle) {
		return ProjectErrInvalidHandle()
	}

	// Every project starts as a draft, whatever the client asked for. Create
	// accepted a status, so a project could be born `active` and never touch
	// the publish flow -- no deployment plan, no record migration, and a
	// namespace that was never swapped into place.
	new.Status = types.ProjectStatusDraft
	new.ArchivedAt = nil

	// FRIA requirement is derived from the deployer answers once, at
	// creation, so the pipeline shape doesn't silently change later.
	new.Config.FriaRequired = friaRequired(new.Config.DeployerCategories)

	if err := svc.uniqueCheck(ctx, new); err != nil {
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

	err := store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
		// Slug uniqueness only matters when one is actually set; empty
		// slugs are allowed to repeat (partial unique index).
		if ns.Slug != "" {
			if existing, e := store.LookupComposeNamespaceBySlug(ctx, s, ns.Slug); e == nil && existing != nil {
				return ProjectErrHandleNotUnique()
			} else if e != nil && !errors.IsNotFound(e) {
				return e
			}
		}

		if err := store.CreateComposeNamespace(ctx, s, ns); err != nil {
			return err
		}

		return store.CreateProject(ctx, s, new)
	})
	if err != nil {
		return err
	}

	if err := label.Create(ctx, svc.store, new); err != nil {
		return err
	}

	// The creator becomes the first member, as a developer: read/write plus
	// the ability to submit the project for publish approval.
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
		if err := store.CreateProjectMember(ctx, svc.store, creator); err != nil {
			return err
		}
	}

	return nil
}

func (svc *project) beforeUpdate(ctx context.Context, upd, res *types.Project) error {
	// Status is not in updateFields, so the generated Update never copies it
	// onto the stored record -- but upd is what gets persisted, so it has to
	// carry the value the store already holds or an update would blank it.
	// Anything a client sent is ignored here, which is the point: a status
	// change means publishing or branching, both of which go through their own
	// endpoints and their own permissions.
	upd.Status = res.Status
	upd.ArchivedAt = res.ArchivedAt

	if upd.Handle != res.Handle {
		if err := svc.uniqueCheck(ctx, upd); err != nil {
			return err
		}
	}

	upd.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()

	// Namespace binding, deployer answers, and FRIA derivation are immutable
	// after creation.
	upd.Config.NamespaceID = res.Config.NamespaceID
	upd.Config.DeployerCategories = res.Config.DeployerCategories
	upd.Config.FriaRequired = res.Config.FriaRequired

	return nil
}

func (svc *project) onDelete(ctx context.Context, s store.Storer, p *types.Project, _ *projectActionProps) error {
	if !svc.ac.CanDeleteProject(ctx, p) {
		return ProjectErrNotAllowedToDelete()
	}

	p.DeletedAt = now()
	p.DeletedBy = a.GetIdentityFromContext(ctx).Identity()

	// The project's namespace follows its lifecycle — leaving it active
	// would squat the slug and block future creations.
	if err := store.UpdateProject(ctx, s, p); err != nil {
		return err
	}
	return svc.setNamespaceDeleted(ctx, s, p.Config.NamespaceID, p.DeletedAt)
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

func (svc *project) onUndelete(ctx context.Context, s store.Storer, p *types.Project, _ *projectActionProps) error {
	if !svc.ac.CanDeleteProject(ctx, p) {
		return ProjectErrNotAllowedToDelete()
	}

	p.DeletedAt = nil
	p.DeletedBy = 0

	if err := store.UpdateProject(ctx, s, p); err != nil {
		return err
	}
	return svc.setNamespaceDeleted(ctx, s, p.Config.NamespaceID, nil)
}

// Archive puts a project on the shelf without touching its lifecycle status.
//
// Archiving used to be a status value, which forced the generic update to
// accept a status from the client -- and that is what let a live project be
// flipped back to "draft" and have its schema edited underneath its own
// records. Keeping the two apart also makes coming back off the shelf
// answerable: status was never overwritten, so Unarchive has nothing to guess.
//
// Gated on update rather than delete: shelving is an edit to how the project is
// presented, not a removal, and it is reversible by anyone who could do it.
func (svc *project) Archive(ctx context.Context, projectID uint64) (*types.Project, error) {
	return svc.setArchived(ctx, projectID, now())
}

// Unarchive brings an archived project back into the working set.
func (svc *project) Unarchive(ctx context.Context, projectID uint64) (*types.Project, error) {
	return svc.setArchived(ctx, projectID, nil)
}

func (svc *project) setArchived(ctx context.Context, projectID uint64, at *time.Time) (*types.Project, error) {
	p, err := loadProject(ctx, svc.store, projectID)
	if err != nil {
		return nil, err
	}

	if !svc.ac.CanUpdateProject(ctx, p) {
		return nil, ProjectErrNotAllowedToUpdate()
	}

	// Idempotent on purpose: archiving an archived project is a no-op rather
	// than an error, so a double-click or a retried request cannot fail.
	if (p.ArchivedAt == nil) == (at == nil) {
		return p, nil
	}

	p.ArchivedAt = at
	p.UpdatedAt = now()
	p.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()

	if err = store.UpdateProject(ctx, svc.store, p); err != nil {
		return nil, err
	}

	return p, nil
}

// --- members ---
//
// Membership is CHAIN-wide, not per revision (ruled 2026-07-29): every read and
// write below resolves the project it was handed to that project's chain ROOT,
// so one list serves every revision. Callers keep passing whichever revision
// they have open and need to know nothing about this.
//
// Before, member rows were written against the revision they were added on, and
// CreateRevision copied none of them — so a freshly branched draft had zero
// members, which resolves to the fallback role preset: no write, no approval,
// no publish. Every branched revision was unusable by everyone, including the
// person who branched it.
//
// The access check still runs against the project as ASKED FOR, not the root:
// permission to read or manage a revision is a question about that revision.

func (svc *project) onSearchMembers(ctx context.Context, _ *projectActionProps, filter types.ProjectMemberFilter) (set types.ProjectMemberSet, f types.ProjectMemberFilter, err error) {
	var p *types.Project
	if p, err = loadProject(ctx, svc.store, filter.ProjectID); err != nil {
		return
	}

	if !svc.ac.CanReadProject(ctx, p) {
		return set, f, ProjectErrNotAllowedToRead()
	}

	filter.ProjectID = p.RootProjectID()

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

	// Written against the chain root, so the membership is the project's rather
	// than this revision's.
	m.ProjectID = p.RootProjectID()

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
	if existing, err = store.LookupProjectMemberByProjectIDUserID(ctx, svc.store, p.RootProjectID(), m.UserID); errors.IsNotFound(err) {
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
	if existing, err = store.LookupProjectMemberByProjectIDUserID(ctx, svc.store, p.RootProjectID(), userID); errors.IsNotFound(err) {
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

// validateProjectStatus is gone with the client-supplied status it existed to
// check: create always writes draft, update never writes status at all, and
// every other transition is made by CreateRevision or Publish with a constant.
// Membership in the enum was never the interesting question anyway -- an
// enum-valid "active" on a project that had never been published was exactly
// the bug.

// friaRequired: any positive deployer-category answer makes a Fundamental
// Rights Impact Assessment step mandatory (EU AI Act Art. 27).
func friaRequired(d types.ProjectDeployerCategories) bool {
	return d.PublicAuthorityAnnex3 || d.PrivateEssentialServices || d.InsuranceBanking
}
