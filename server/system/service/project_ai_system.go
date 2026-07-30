package service

import (
	"context"

	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/handle"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

// The CRUD method bodies (FindByID, Search, Create, Update, DeleteByID and
// loadProjectAiSystem) are generated in project_ai_system.gen.go from
// system/project_ai_system.cue.
//
// This file owns the struct, access-controller interface, constructor, the
// public service contract, entry (resource member) management and the
// before-create / before-update hooks the generated Create / Update call into.

type (
	projectAiSystemAccessController interface {
		CanCreateProjectAiSystem(ctx context.Context) bool
		CanSearchProjectAiSystems(ctx context.Context) bool
		CanReadProjectAiSystem(ctx context.Context, r *types.ProjectAiSystem) bool
		CanUpdateProjectAiSystem(ctx context.Context, r *types.ProjectAiSystem) bool
		CanDeleteProjectAiSystem(ctx context.Context, r *types.ProjectAiSystem) bool
		CanManageResourcesOnProjectAiSystem(ctx context.Context, r *types.ProjectAiSystem) bool
	}

	ProjectAiSystemService interface {
		FindByID(ctx context.Context, id uint64) (*types.ProjectAiSystem, error)
		Search(ctx context.Context, f types.ProjectAiSystemFilter) (types.ProjectAiSystemSet, types.ProjectAiSystemFilter, error)
		Create(ctx context.Context, g *types.ProjectAiSystem) (*types.ProjectAiSystem, error)
		Update(ctx context.Context, g *types.ProjectAiSystem) (*types.ProjectAiSystem, error)
		DeleteByID(ctx context.Context, id uint64) error

		MemberList(ctx context.Context, projectAiSystemID uint64) (types.ProjectAiSystemEntrySet, error)
		MemberAdd(ctx context.Context, projectAiSystemID uint64, resourceRef string) error
		MemberRemove(ctx context.Context, projectAiSystemID uint64, resourceRef string) error
	}
)

func ProjectAiSystem() *projectAiSystem {
	return &projectAiSystem{
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		ac:        DefaultAccessControl,
	}
}

// beforeCreate validates the handle, resolves the owning tenant from the
// project and enforces handle uniqueness before the generated Create assigns
// the ID / timestamps and persists.
func (svc *projectAiSystem) beforeCreate(ctx context.Context, new *types.ProjectAiSystem) error {
	if !handle.IsValid(new.Handle) {
		return ProjectAiSystemErrInvalidHandle()
	}

	p, err := loadProject(ctx, svc.store, new.ProjectID)
	if err != nil {
		return err
	}

	if err = svc.uniqueCheck(ctx, new); err != nil {
		return err
	}

	new.TenantID = p.TenantID
	return nil
}

// beforeUpdate merges the mutable fields (handle, risk class, meta) with
// validation and uniqueness checks. The generated Update has already loaded
// `existing`, run the access check and the stale-data guard.
func (svc *projectAiSystem) beforeUpdate(ctx context.Context, upd, existing *types.ProjectAiSystem) error {
	if upd.Handle != "" && !handle.IsValid(upd.Handle) {
		return ProjectAiSystemErrInvalidHandle()
	}

	if upd.Handle != "" && upd.Handle != existing.Handle {
		if err := svc.uniqueCheck(ctx, &types.ProjectAiSystem{ID: upd.ID, ProjectID: existing.ProjectID, Handle: upd.Handle}); err != nil {
			return err
		}
		existing.Handle = upd.Handle
	}

	if upd.RiskClass != "" {
		existing.RiskClass = upd.RiskClass
	}

	if upd.Meta.Short != "" {
		existing.Meta.Short = upd.Meta.Short
	}
	if upd.Meta.Description != "" {
		existing.Meta.Description = upd.Meta.Description
	}
	if upd.Meta.IntendedPurpose != "" {
		existing.Meta.IntendedPurpose = upd.Meta.IntendedPurpose
	}

	return nil
}

// --- entries (resource members) ---

func (svc *projectAiSystem) onMemberList(ctx context.Context, _ *projectAiSystemActionProps, projectAiSystemID uint64) (set types.ProjectAiSystemEntrySet, err error) {
	g, err := loadProjectAiSystem(ctx, svc.store, projectAiSystemID)
	if err != nil {
		return nil, err
	}

	if !svc.ac.CanReadProjectAiSystem(ctx, g) {
		return nil, ProjectAiSystemErrNotAllowedToRead()
	}

	set, _, err = store.SearchProjectAiSystemEntrys(ctx, svc.store, types.ProjectAiSystemEntryFilter{
		ProjectAiSystemID: projectAiSystemID,
	})
	return set, err
}

func (svc *projectAiSystem) onMemberAdd(ctx context.Context, _ *projectAiSystemActionProps, projectAiSystemID uint64, resourceRef string) error {
	g, err := loadProjectAiSystem(ctx, svc.store, projectAiSystemID)
	if err != nil {
		return err
	}

	if !svc.ac.CanManageResourcesOnProjectAiSystem(ctx, g) {
		return ProjectAiSystemErrNotAllowedToManageResources()
	}

	if resourceRef == "" {
		return ProjectAiSystemErrInvalidID()
	}

	if existing, e := store.LookupProjectAiSystemEntryByProjectAiSystemIDResourceRef(ctx, svc.store, projectAiSystemID, resourceRef); e == nil && existing != nil {
		return nil
	} else if e != nil && !errors.IsNotFound(e) {
		return e
	}

	return store.CreateProjectAiSystemEntry(ctx, svc.store, &types.ProjectAiSystemEntry{
		ProjectAiSystemID: projectAiSystemID,
		ResourceRef:       resourceRef,
		CreatedAt:         *now(),
	})
}

func (svc *projectAiSystem) onMemberRemove(ctx context.Context, _ *projectAiSystemActionProps, projectAiSystemID uint64, resourceRef string) error {
	g, err := loadProjectAiSystem(ctx, svc.store, projectAiSystemID)
	if err != nil {
		return err
	}

	if !svc.ac.CanManageResourcesOnProjectAiSystem(ctx, g) {
		return ProjectAiSystemErrNotAllowedToManageResources()
	}

	return store.DeleteProjectAiSystemEntryByProjectAiSystemIDResourceRef(ctx, svc.store, projectAiSystemID, resourceRef)
}

// --- helpers ---

func (svc *projectAiSystem) uniqueCheck(ctx context.Context, g *types.ProjectAiSystem) error {
	if g.Handle == "" {
		return nil
	}
	if e, err := store.LookupProjectAiSystemByProjectIDHandle(ctx, svc.store, g.ProjectID, g.Handle); err == nil && e != nil && e.ID != g.ID {
		return ProjectAiSystemErrHandleNotUnique()
	} else if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}
