package rest

import (
	"context"

	"github.com/crusttech/crust/internal/permissions"
	"github.com/crusttech/crust/system/internal/service"
	"github.com/crusttech/crust/system/rest/request"
)

type (
	Permissions struct {
		ac permissionsAccessController
	}

	permissionsAccessController interface {
		Effective(context.Context) permissions.EffectiveSet
	}
)

func (Permissions) New() *Permissions {
	return &Permissions{
		ac: service.DefaultAccessControl,
	}
}

func (ctrl *Permissions) Effective(ctx context.Context, r *request.PermissionsEffective) (interface{}, error) {
	return ctrl.ac.Effective(ctx), nil
}

func (ctrl *Permissions) List(ctx context.Context, r *request.PermissionsList) (interface{}, error) {
	return "not implemented", nil
	// return ctrl.svc.rules.With(ctx).List()
}

func (ctrl *Permissions) Read(ctx context.Context, r *request.PermissionsRead) (interface{}, error) {
	return "not implemented", nil
	// return ctrl.svc.rules.With(ctx).Read(r.RoleID)
}

func (ctrl *Permissions) Delete(ctx context.Context, r *request.PermissionsDelete) (interface{}, error) {
	return "not implemented", nil
	// return ctrl.svc.rules.With(ctx).Delete(r.RoleID)
}

func (ctrl *Permissions) Update(ctx context.Context, r *request.PermissionsUpdate) (interface{}, error) {
	return "not implemented", nil
	// return ctrl.svc.rules.With(ctx).Update(r.RoleID, r.Permissions)
}
