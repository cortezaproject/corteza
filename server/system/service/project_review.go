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
	if new.Title == "" {
		new.Title = firstNonEmpty(new.Description, "(untitled)")
	}
	if new.Status == "" {
		new.Status = "Open"
	}
	new.CreatedBy = a.GetIdentityFromContext(ctx).Identity()
	return nil
}

func (svc *projectReview) beforeUpdate(ctx context.Context, upd, _ *types.ProjectReview) error {
	upd.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()
	return nil
}
