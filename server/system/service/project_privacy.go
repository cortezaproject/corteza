package service

import (
	"context"

	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/system/types"
)

func ProjectPrivacy() *projectPrivacy {
	return &projectPrivacy{
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
	}
}

func (svc *projectPrivacy) beforeCreate(ctx context.Context, new *types.ProjectPrivacy) error {
	if new.ProjectID == 0 {
		return ProjectPrivacyErrMissingProject()
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

func (svc *projectPrivacy) beforeUpdate(ctx context.Context, upd, _ *types.ProjectPrivacy) error {
	upd.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()
	return nil
}

func (svc *projectPrivacy) beforeDelete(ctx context.Context, res *types.ProjectPrivacy) error {
	res.DeletedBy = a.GetIdentityFromContext(ctx).Identity()
	return nil
}
