package rest

import (
	"context"

	"github.com/crusttech/crust/sam/rest/request"
	"github.com/crusttech/crust/sam/service"
	_ "github.com/crusttech/crust/sam/types"
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

func (ctrl *Permissions) GetTeam(ctx context.Context, r *request.PermissionsGetTeam) (interface{}, error) {
	return ctrl.svc.perms.Get(r.Team, r.Scope, r.Resource)
}

func (ctrl *Permissions) SetTeam(ctx context.Context, r *request.PermissionsSetTeam) (interface{}, error) {
	return ctrl.svc.perms.Set(r.Team, r.Permissions)
}
