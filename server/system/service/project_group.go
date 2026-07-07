package service

import (
	"context"

	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/handle"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

// The CRUD method bodies (FindByID, Search, Create, Update, DeleteByID and
// loadProjectGroup) are generated in project_group.gen.go from
// system/project_group.cue.
//
// This file owns the struct, access-controller interface, constructor, the
// public service contract, member management and the before-create /
// before-update hooks the generated Create / Update call into.

type (
	projectGroupAccessController interface {
		CanCreateProjectGroup(ctx context.Context) bool
		CanSearchProjectGroups(ctx context.Context) bool
		CanReadProjectGroup(ctx context.Context, r *types.ProjectGroup) bool
		CanUpdateProjectGroup(ctx context.Context, r *types.ProjectGroup) bool
		CanDeleteProjectGroup(ctx context.Context, r *types.ProjectGroup) bool
		CanManageMembersOnProjectGroup(ctx context.Context, r *types.ProjectGroup) bool
	}

	ProjectGroupService interface {
		FindByID(ctx context.Context, id uint64) (*types.ProjectGroup, error)
		Search(ctx context.Context, f types.ProjectGroupFilter) (types.ProjectGroupSet, types.ProjectGroupFilter, error)
		Create(ctx context.Context, g *types.ProjectGroup) (*types.ProjectGroup, error)
		Update(ctx context.Context, g *types.ProjectGroup) (*types.ProjectGroup, error)
		DeleteByID(ctx context.Context, id uint64) error

		MemberList(ctx context.Context, projectGroupID uint64) (types.ProjectGroupEntrySet, error)
		MemberAdd(ctx context.Context, projectGroupID uint64, resourceRef string) error
		MemberRemove(ctx context.Context, projectGroupID uint64, resourceRef string) error
	}
)

// beforeCreate validates the handle, resolves the owning tenant from the
// project and enforces handle uniqueness before the generated Create assigns
// the ID / timestamps and persists.
func (svc *projectGroup) beforeCreate(ctx context.Context, new *types.ProjectGroup) error {
	if !handle.IsValid(new.Handle) {
		return ProjectGroupErrInvalidHandle()
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

// beforeUpdate merges the mutable fields (handle, meta) with validation and
// uniqueness checks. The generated Update has already loaded `existing`, run
// the access check and the stale-data guard.
func (svc *projectGroup) beforeUpdate(ctx context.Context, upd, existing *types.ProjectGroup) error {
	if upd.Handle != "" && !handle.IsValid(upd.Handle) {
		return ProjectGroupErrInvalidHandle()
	}

	if upd.Handle != "" && upd.Handle != existing.Handle {
		if err := svc.uniqueCheck(ctx, &types.ProjectGroup{ID: upd.ID, ProjectID: existing.ProjectID, Handle: upd.Handle}); err != nil {
			return err
		}
		existing.Handle = upd.Handle
	}

	if upd.Meta.Short != "" {
		existing.Meta.Short = upd.Meta.Short
	}
	if upd.Meta.Description != "" {
		existing.Meta.Description = upd.Meta.Description
	}

	return nil
}

// --- members ---

func (svc *projectGroup) onMemberList(ctx context.Context, _ *projectGroupActionProps, projectGroupID uint64) (set types.ProjectGroupEntrySet, err error) {
	g, err := loadProjectGroup(ctx, svc.store, projectGroupID)
	if err != nil {
		return nil, err
	}

	if !svc.ac.CanReadProjectGroup(ctx, g) {
		return nil, ProjectGroupErrNotAllowedToRead()
	}

	set, _, err = store.SearchProjectGroupEntrys(ctx, svc.store, types.ProjectGroupEntryFilter{
		ProjectGroupID: projectGroupID,
	})
	return set, err
}

func (svc *projectGroup) onMemberAdd(ctx context.Context, _ *projectGroupActionProps, projectGroupID uint64, resourceRef string) error {
	g, err := loadProjectGroup(ctx, svc.store, projectGroupID)
	if err != nil {
		return err
	}

	if !svc.ac.CanManageMembersOnProjectGroup(ctx, g) {
		return ProjectGroupErrNotAllowedToManageMembers()
	}

	if resourceRef == "" {
		return ProjectGroupErrInvalidID()
	}

	if existing, e := store.LookupProjectGroupEntryByProjectGroupIDResourceRef(ctx, svc.store, projectGroupID, resourceRef); e == nil && existing != nil {
		return nil
	} else if e != nil && !errors.IsNotFound(e) {
		return e
	}

	return store.CreateProjectGroupEntry(ctx, svc.store, &types.ProjectGroupEntry{
		ProjectGroupID: projectGroupID,
		ResourceRef:    resourceRef,
		CreatedAt:      *now(),
	})
}

func (svc *projectGroup) onMemberRemove(ctx context.Context, _ *projectGroupActionProps, projectGroupID uint64, resourceRef string) error {
	g, err := loadProjectGroup(ctx, svc.store, projectGroupID)
	if err != nil {
		return err
	}

	if !svc.ac.CanManageMembersOnProjectGroup(ctx, g) {
		return ProjectGroupErrNotAllowedToManageMembers()
	}

	return store.DeleteProjectGroupEntryByProjectGroupIDResourceRef(ctx, svc.store, projectGroupID, resourceRef)
}

// --- helpers ---

func (svc *projectGroup) uniqueCheck(ctx context.Context, g *types.ProjectGroup) error {
	if g.Handle == "" {
		return nil
	}
	if e, err := store.LookupProjectGroupByProjectIDHandle(ctx, svc.store, g.ProjectID, g.Handle); err == nil && e != nil && e.ID != g.ID {
		return ProjectGroupErrHandleNotUnique()
	} else if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}
