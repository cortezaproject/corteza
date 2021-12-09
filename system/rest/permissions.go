package rest

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/system/rest/request"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	Permissions struct {
		ac permissionsAccessController
	}

	permissionsAccessController interface {
		Effective(context.Context, ...rbac.Resource) rbac.EffectiveSet
		List() []map[string]string
		FindRulesByRoleID(context.Context, uint64) (rbac.RuleSet, error)
		CloneRulesByRoleID(ctx context.Context, roleID uint64, toRoleID ...uint64) error
		Grant(ctx context.Context, rr ...*rbac.Rule) error
	}
)

func (Permissions) New() *Permissions {
	return &Permissions{
		ac: service.DefaultAccessControl,
	}
}

func (ctrl Permissions) Effective(ctx context.Context, r *request.PermissionsEffective) (interface{}, error) {
	return ctrl.ac.Effective(ctx, types.Component{}), nil
}

func (ctrl Permissions) List(ctx context.Context, r *request.PermissionsList) (interface{}, error) {
	return ctrl.ac.List(), nil
}

func (ctrl Permissions) Read(ctx context.Context, r *request.PermissionsRead) (interface{}, error) {
	return ctrl.ac.FindRulesByRoleID(ctx, r.RoleID)
}

func (ctrl Permissions) Delete(ctx context.Context, r *request.PermissionsDelete) (interface{}, error) {
	rr, err := ctrl.ac.FindRulesByRoleID(ctx, r.RoleID)
	if err != nil {
		return nil, err
	}

	for _, r := range rr {
		r.Access = rbac.Inherit
	}

	return api.OK(), ctrl.ac.Grant(ctx, rr...)
}

func (ctrl Permissions) Update(ctx context.Context, r *request.PermissionsUpdate) (interface{}, error) {
	for _, rule := range r.Rules {
		// Make sure everything is properly set
		rule.RoleID = r.RoleID
	}

	return api.OK(), ctrl.ac.Grant(ctx, r.Rules...)
}

func (ctrl Permissions) Clone(ctx context.Context, r *request.PermissionsClone) (interface{}, error) {
	// Clone rules from role S to role T
	return api.OK(), ctrl.ac.CloneRulesByRoleID(ctx, r.RoleID, payload.ParseUint64s(r.CloneToRoleID)...)
}
