package rest

import (
	"context"

	"github.com/cortezaproject/corteza-server/compose/rest/request"
	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
)

type (
	Permissions struct {
		ac permissionsAccessController
	}

	rbacResWrap struct {
		res string
	}

	permissionsAccessController interface {
		Effective(context.Context, ...rbac.Resource) rbac.EffectiveSet
		Evaluate(ctx context.Context, user uint64, roles []uint64, rr ...rbac.Resource) (ee rbac.EvaluatedSet, err error)
		List() []map[string]string
		FindRulesByRoleID(context.Context, uint64) (rbac.RuleSet, error)
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

func (ctrl Permissions) Evaluate(ctx context.Context, r *request.PermissionsEvaluate) (interface{}, error) {
	in := make([]rbac.Resource, 0, len(r.Resource))
	for _, res := range r.Resource {
		in = append(in, rbacResWrap{res: res})
	}

	return ctrl.ac.Evaluate(ctx, r.UserID, r.RoleID, in...)
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

func (ar rbacResWrap) RbacResource() string {
	return ar.res
}
