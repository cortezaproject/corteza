package rest

import (
	"context"

	"github.com/crusttech/crust/internal/permissions"
	"github.com/crusttech/crust/messaging/internal/service"
	"github.com/crusttech/crust/messaging/rest/request"
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
