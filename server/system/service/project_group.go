package service

import (
	"context"

	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/handle"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

type (
	projectGroup struct {
		actionlog actionlog.Recorder
		store     store.Storer
		ac        projectGroupAccessController
	}

	projectGroupAccessController interface {
		CanCreateProjectGroup(ctx context.Context) bool
		CanSearchProjectGroups(ctx context.Context) bool
		CanReadProjectGroup(ctx context.Context, g *types.ProjectGroup) bool
		CanUpdateProjectGroup(ctx context.Context, g *types.ProjectGroup) bool
		CanDeleteProjectGroup(ctx context.Context, g *types.ProjectGroup) bool
		CanManageMembersOnProjectGroup(ctx context.Context, g *types.ProjectGroup) bool
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

func ProjectGroup() *projectGroup {
	return &projectGroup{
		actionlog: DefaultActionlog,
		ac:        DefaultAccessControl,
		store:     DefaultStore,
	}
}

func (svc *projectGroup) FindByID(ctx context.Context, id uint64) (g *types.ProjectGroup, err error) {
	if g, err = loadProjectGroup(ctx, svc.store, id); err != nil {
		return nil, err
	}

	if !svc.ac.CanReadProjectGroup(ctx, g) {
		return nil, ProjectGroupErrNotAllowedToRead()
	}

	return g, nil
}

func (svc *projectGroup) Search(ctx context.Context, f types.ProjectGroupFilter) (set types.ProjectGroupSet, rf types.ProjectGroupFilter, err error) {
	if !svc.ac.CanSearchProjectGroups(ctx) {
		return nil, f, ProjectGroupErrNotAllowedToSearch()
	}

	f.Check = func(res *types.ProjectGroup) (bool, error) {
		if !svc.ac.CanReadProjectGroup(ctx, res) {
			return false, nil
		}
		return true, nil
	}

	return store.SearchProjectGroups(ctx, svc.store, f)
}

func (svc *projectGroup) Create(ctx context.Context, new *types.ProjectGroup) (g *types.ProjectGroup, err error) {
	if !svc.ac.CanCreateProjectGroup(ctx) {
		return nil, ProjectGroupErrNotAllowedToCreate()
	}

	if !handle.IsValid(new.Handle) {
		return nil, ProjectGroupErrInvalidHandle()
	}

	p, err := loadProject(ctx, svc.store, new.ProjectID)
	if err != nil {
		return nil, err
	}

	if err = svc.uniqueCheck(ctx, new); err != nil {
		return nil, err
	}

	new.ID = nextID()
	new.TenantID = p.TenantID
	new.CreatedAt = *now()

	if err = store.CreateProjectGroup(ctx, svc.store, new); err != nil {
		return nil, err
	}

	return new, nil
}

func (svc *projectGroup) Update(ctx context.Context, upd *types.ProjectGroup) (g *types.ProjectGroup, err error) {
	existing, err := loadProjectGroup(ctx, svc.store, upd.ID)
	if err != nil {
		return nil, err
	}

	if !svc.ac.CanUpdateProjectGroup(ctx, existing) {
		return nil, ProjectGroupErrNotAllowedToUpdate()
	}

	if isStale(upd.UpdatedAt, existing.UpdatedAt, existing.CreatedAt) {
		return nil, ProjectGroupErrStaleData()
	}

	if upd.Handle != "" && !handle.IsValid(upd.Handle) {
		return nil, ProjectGroupErrInvalidHandle()
	}

	if upd.Handle != "" && upd.Handle != existing.Handle {
		if err = svc.uniqueCheck(ctx, &types.ProjectGroup{ID: upd.ID, ProjectID: existing.ProjectID, Handle: upd.Handle}); err != nil {
			return nil, err
		}
		existing.Handle = upd.Handle
	}

	if upd.Meta.Short != "" {
		existing.Meta.Short = upd.Meta.Short
	}
	if upd.Meta.Description != "" {
		existing.Meta.Description = upd.Meta.Description
	}

	existing.UpdatedAt = now()

	if err = store.UpdateProjectGroup(ctx, svc.store, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (svc *projectGroup) DeleteByID(ctx context.Context, id uint64) (err error) {
	g, err := loadProjectGroup(ctx, svc.store, id)
	if err != nil {
		return err
	}

	if !svc.ac.CanDeleteProjectGroup(ctx, g) {
		return ProjectGroupErrNotAllowedToDelete()
	}

	g.DeletedAt = now()
	return store.UpdateProjectGroup(ctx, svc.store, g)
}

// --- members ---

func (svc *projectGroup) MemberList(ctx context.Context, projectGroupID uint64) (types.ProjectGroupEntrySet, error) {
	g, err := loadProjectGroup(ctx, svc.store, projectGroupID)
	if err != nil {
		return nil, err
	}

	if !svc.ac.CanReadProjectGroup(ctx, g) {
		return nil, ProjectGroupErrNotAllowedToRead()
	}

	set, _, err := store.SearchProjectGroupEntrys(ctx, svc.store, types.ProjectGroupEntryFilter{
		ProjectGroupID: projectGroupID,
	})
	return set, err
}

func (svc *projectGroup) MemberAdd(ctx context.Context, projectGroupID uint64, resourceRef string) error {
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

func (svc *projectGroup) MemberRemove(ctx context.Context, projectGroupID uint64, resourceRef string) error {
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

func loadProjectGroup(ctx context.Context, s store.ProjectGroups, id uint64) (res *types.ProjectGroup, err error) {
	if id == 0 {
		return nil, ProjectGroupErrInvalidID()
	}

	if res, err = store.LookupProjectGroupByID(ctx, s, id); errors.IsNotFound(err) {
		return nil, ProjectGroupErrNotFound()
	}

	return
}
