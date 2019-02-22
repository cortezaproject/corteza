package rest

import (
	"context"

	"github.com/pkg/errors"

	"github.com/crusttech/crust/system/rest/request"
	"github.com/crusttech/crust/system/service"
)

var _ = errors.Wrap

type (
	Permission struct {
		svc struct {
			perm service.PermissionService
		}
	}
)

func (Permission) New() *Permission {
	ctrl := &Permission{}
	ctrl.svc.perm = service.DefaultPermission
	return ctrl
}

func (ctrl *Permission) Read(ctx context.Context, r *request.PermissionRead) (interface{}, error) {
	return ctrl.svc.perm.Read(r.RoleID)
}

func (ctrl *Permission) Delete(ctx context.Context, r *request.PermissionDelete) (interface{}, error) {
	return ctrl.svc.perm.Delete(r.RoleID)
}

func (ctrl *Permission) Update(ctx context.Context, r *request.PermissionUpdate) (interface{}, error) {
	return ctrl.svc.perm.Update(r.RoleID, r.Permissions)
}
