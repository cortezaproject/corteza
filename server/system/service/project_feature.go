package service

import (
	"context"

	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/system/types"
)

// The struct, access-controller interface, and CRUD method bodies (FindByID,
// Search, Create, Update, DeleteByID, loadProjectFeature) are generated in
// project_feature.gen.go from the CUE model.
//
// This file owns the constructor and the before-create / before-update hooks
// the generated Create / Update call into.

func ProjectFeature() *projectFeature {
	return &projectFeature{
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
	}
}

func (svc *projectFeature) beforeCreate(ctx context.Context, new *types.ProjectFeature) error {
	if new.Title == "" {
		new.Title = firstNonEmpty(new.Description, "(untitled)")
	}
	if new.Status == "" {
		new.Status = "Open"
	}
	new.CreatedBy = a.GetIdentityFromContext(ctx).Identity()
	return nil
}

func (svc *projectFeature) beforeUpdate(ctx context.Context, upd, _ *types.ProjectFeature) error {
	upd.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()
	return nil
}
