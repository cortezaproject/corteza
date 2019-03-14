package rest

import (
	"context"

	"github.com/pkg/errors"

	"github.com/crusttech/crust/system/rest/request"
	"github.com/crusttech/crust/system/internal/service"
)

var _ = errors.Wrap

type (
	Permissions struct {
		svc struct {
			rules service.RulesService
			perm  service.PermissionsService
		}
	}
)

func (Permissions) New() *Permissions {
	ctrl := &Permissions{}
	ctrl.svc.rules = service.DefaultRules
	ctrl.svc.perm = service.DefaultPermissions
	return ctrl
}

func (ctrl *Permissions) List(ctx context.Context, r *request.PermissionsList) (interface{}, error) {
	return ctrl.svc.rules.With(ctx).List()
}

func (ctrl *Permissions) Effective(ctx context.Context, r *request.PermissionsEffective) (interface{}, error) {
	return ctrl.svc.perm.With(ctx).Effective()
}

func (ctrl *Permissions) Read(ctx context.Context, r *request.PermissionsRead) (interface{}, error) {
	return ctrl.svc.rules.With(ctx).Read(r.RoleID)
}

func (ctrl *Permissions) Delete(ctx context.Context, r *request.PermissionsDelete) (interface{}, error) {
	return ctrl.svc.rules.With(ctx).Delete(r.RoleID)
}

func (ctrl *Permissions) Update(ctx context.Context, r *request.PermissionsUpdate) (interface{}, error) {
	return ctrl.svc.rules.With(ctx).Update(r.RoleID, r.Permissions)
}
