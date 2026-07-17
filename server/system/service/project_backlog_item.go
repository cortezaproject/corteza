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

func (svc *projectBacklogItem) beforeCreate(ctx context.Context, new *types.ProjectBacklogItem) error {
	if new.ProjectID == 0 {
		return ProjectBacklogItemErrMissingProject()
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

func (svc *projectBacklogItem) beforeUpdate(ctx context.Context, upd, _ *types.ProjectBacklogItem) error {
	if upd.Category != "" && !projectBacklogItemValidCategories[upd.Category] {
		return ProjectBacklogItemErrInvalidCategory()
	}
	upd.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()
	return nil
}

func (svc *projectBacklogItem) beforeDelete(ctx context.Context, res *types.ProjectBacklogItem) error {
	res.DeletedBy = a.GetIdentityFromContext(ctx).Identity()
	return nil
}
