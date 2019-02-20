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
			perm service.PermissionService
		}
	}
)

func (Permissions) New() *Permissions {
	return &Permissions{}
}

func (ctrl *Permissions) Get(ctx context.Context, r *request.PermissionsGet) (interface{}, error) {
	return ctrl.svc.perm.Get(r.RoleID)
}

func (ctrl *Permissions) Delete(ctx context.Context, r *request.PermissionsDelete) (interface{}, error) {
	return ctrl.svc.perm.Delete(r.RoleID)
}

func (ctrl *Permissions) Update(ctx context.Context, r *request.PermissionsUpdate) (interface{}, error) {
	return ctrl.svc.perm.Update(r.RoleID, r.Permissions)
}
