package rest

import (
	"context"

	"github.com/crusttech/crust/crm/rest/request"
	"github.com/crusttech/crust/crm/service"

	"github.com/pkg/errors"
)

var _ = errors.Wrap

type Permissions struct {
	svc struct {
		perm service.PermissionsService
	}
}

func (Permissions) New() *Permissions {
	ctrl := &Permissions{}
	ctrl.svc.perm = service.DefaultPermissions
	return ctrl
}

func (ctrl *Permissions) Effective(ctx context.Context, r *request.PermissionsEffective) (interface{}, error) {
	return ctrl.svc.perm.With(ctx).Effective()
}
