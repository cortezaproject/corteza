package rest

import (
	"context"

	"github.com/crusttech/crust/messaging/rest/request"
	"github.com/crusttech/crust/messaging/service"
	_ "github.com/crusttech/crust/messaging/types"
)

type Permissions struct {
	svc struct {
		perms service.PermissionsService
	}
}

func (Permissions) New() *Permissions {
	ctrl := &Permissions{}
	ctrl.svc.perms = service.DefaultPermissions
	return ctrl
}

func (ctrl *Permissions) List(ctx context.Context, r *request.PermissionsList) (interface{}, error) {
	return ctrl.svc.perms.List()
}

func (ctrl *Permissions) Get(ctx context.Context, r *request.PermissionsGet) (interface{}, error) {
	return ctrl.svc.perms.Get(r.TeamID, r.Resource)
}

func (ctrl *Permissions) Set(ctx context.Context, r *request.PermissionsSet) (interface{}, error) {
	return ctrl.svc.perms.Set(r.TeamID, r.Permissions)
}

func (ctrl *Permissions) Scopes(ctx context.Context, r *request.PermissionsScopes) (interface{}, error) {
	return ctrl.svc.perms.Scopes(r.Scope)
}
