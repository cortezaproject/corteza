package service

import (
	"context"

	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/system/types"
)

// projectBacklogItemValidCategories are the 5 category tables a backlog item
// (sub-issue) can be filed against; Category + EventID together form the
// polymorphic reference to the parent category item.
var projectBacklogItemValidCategories = map[string]bool{
	"incident": true,
	"task":     true,
	"feature":  true,
	"privacy":  true,
	"review":   true,
}

func ProjectBacklogItem() *projectBacklogItem {
	return &projectBacklogItem{
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
	}
}

// beforeSearch normalises the project scope to the chain root, mirroring what
// beforeCreate already does on the way in. See rootProjectID (project.go).
func (svc *projectBacklogItem) beforeSearch(ctx context.Context, f *types.ProjectBacklogItemFilter) error {
	f.ProjectID = rootProjectID(ctx, svc.store, f.ProjectID)
	return nil
}

func (svc *projectBacklogItem) beforeCreate(ctx context.Context, new *types.ProjectBacklogItem) error {
	if new.ProjectID == 0 {
		return ProjectBacklogItemErrMissingProject()
	}

	// Work items always file against the chain root, never a revision, so a
	// revision reads like a milestone over one shared item pool. New items are
	// deliberately left unassigned (RevisionID stays 0) here.
	project, err := loadProject(ctx, svc.store, new.ProjectID)
	if err != nil {
		return err
	}
	new.ProjectID = project.RootProjectID()

	if ok, err := revisionInChain(ctx, svc.store, new.ProjectID, new.RevisionID); err != nil {
		return err
	} else if !ok {
		return ProjectBacklogItemErrInvalidRevision()
	}

	if new.EventID == 0 {
		return ProjectBacklogItemErrMissingEvent()
	}
	if !projectBacklogItemValidCategories[new.Category] {
		return ProjectBacklogItemErrInvalidCategory()
	}
	if new.Title == "" {
		new.Title = firstNonEmpty(new.Description, "(untitled)")
	}
	if new.Status == "" {
		new.Status = "Open"
	}
	new.CreatedBy = a.GetIdentityFromContext(ctx).Identity()
	return nil
}

func (svc *projectBacklogItem) beforeUpdate(ctx context.Context, upd, res *types.ProjectBacklogItem) error {
	if upd.Category != "" && !projectBacklogItemValidCategories[upd.Category] {
		return ProjectBacklogItemErrInvalidCategory()
	}
	if ok, err := revisionInChain(ctx, svc.store, res.ProjectID, upd.RevisionID); err != nil {
		return err
	} else if !ok {
		return ProjectBacklogItemErrInvalidRevision()
	}
	upd.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()
	return nil
}

func (svc *projectBacklogItem) beforeDelete(ctx context.Context, res *types.ProjectBacklogItem) error {
	res.DeletedBy = a.GetIdentityFromContext(ctx).Identity()
	return nil
}
