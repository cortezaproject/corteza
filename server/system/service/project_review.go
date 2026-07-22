package service

import (
	"context"

	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/system/types"
)

func ProjectReview() *projectReview {
	return &projectReview{
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
	}
}

func (svc *projectReview) beforeCreate(ctx context.Context, new *types.ProjectReview) error {
	if new.ProjectID == 0 {
		return ProjectReviewErrMissingProject()
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

func (svc *projectReview) beforeUpdate(ctx context.Context, upd, res *types.ProjectReview) error {
	// Stamp res (the persisted record); updated_by is not in the generated
	// field-copy, so stamping upd would be dropped.
	res.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()
	return nil
}

func (svc *projectReview) beforeDelete(ctx context.Context, res *types.ProjectReview) error {
	res.DeletedBy = a.GetIdentityFromContext(ctx).Identity()
	return nil
}
