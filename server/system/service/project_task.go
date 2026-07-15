package service

import (
	"context"

	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/system/types"
)

func ProjectTask() *projectTask {
	return &projectTask{
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
	}
}

func (svc *projectTask) beforeCreate(ctx context.Context, new *types.ProjectTask) error {
	if new.Title == "" {
		new.Title = firstNonEmpty(new.Description, "(untitled)")
	}
	if new.Status == "" {
		new.Status = "Open"
	}
	new.CreatedBy = a.GetIdentityFromContext(ctx).Identity()
	return nil
}

func (svc *projectTask) beforeUpdate(ctx context.Context, upd, _ *types.ProjectTask) error {
	upd.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()
	return nil
}
