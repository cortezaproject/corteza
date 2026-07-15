package service

import (
	"context"

	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/system/types"
)

func ProjectIncident() *projectIncident {
	return &projectIncident{
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
	}
}

func (svc *projectIncident) beforeCreate(ctx context.Context, new *types.ProjectIncident) error {
	if new.Title == "" {
		new.Title = firstNonEmpty(new.Description, "(untitled)")
	}
	if new.Status == "" {
		new.Status = "Open"
	}
	new.CreatedBy = a.GetIdentityFromContext(ctx).Identity()
	return nil
}

func (svc *projectIncident) beforeUpdate(ctx context.Context, upd, _ *types.ProjectIncident) error {
	upd.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()
	return nil
}
