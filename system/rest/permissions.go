package rest

import (
	"context"

	"github.com/pkg/errors"

	"github.com/crusttech/crust/system/rest/request"
	"github.com/crusttech/crust/system/service"
)

var _ = errors.Wrap

type (
	Permissions struct {
		svc struct {
			rules service.RulesService
		}
	}
)

func (Permissions) New() *Permissions {
	ctrl := &Permissions{}
	ctrl.svc.rules = service.DefaultRules
	return ctrl
}

func (ctrl *Permissions) List(ctx context.Context, r *request.PermissionsList) (interface{}, error) {
	return ctrl.svc.rules.List()
}

func (ctrl *Permissions) Read(ctx context.Context, r *request.PermissionsRead) (interface{}, error) {
	return ctrl.svc.rules.Read(r.RoleID)
}

func (ctrl *Permissions) Delete(ctx context.Context, r *request.PermissionsDelete) (interface{}, error) {
	return ctrl.svc.rules.Delete(r.RoleID)
}

func (ctrl *Permissions) Update(ctx context.Context, r *request.PermissionsUpdate) (interface{}, error) {
	return ctrl.svc.rules.Update(r.RoleID, r.Permissions)
}
