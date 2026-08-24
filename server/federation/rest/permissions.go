package rest

import (
	"context"
	"github.com/crusttech/human/server/federation/rest/request"
	"github.com/crusttech/human/server/federation/service"
	"github.com/crusttech/human/server/federation/types"
	"github.com/crusttech/human/server/pkg/api"
	"github.com/crusttech/human/server/pkg/rbac"
)

type (
	Permissions struct {
		ac permissionsAccessController
	}

	permissionsAccessController interface {
		Effective(context.Context, ...rbac.Resource) rbac.EffectiveSet
		EffectiveFor(context.Context, string) (rbac.EffectiveSet, error)
		Trace(context.Context, uint64, []uint64, ...string) ([]*rbac.Trace, error)
		List() []map[string]string
		FindRulesByRoleID(context.Context, uint64) (rbac.RuleSet, error)
		FindRules(ctx context.Context, roleID uint64, rr ...string) (rbac.RuleSet, error)
		Grant(ctx context.Context, rr ...*rbac.Rule) error
	}
)

func (Permissions) New() *Permissions {
	return &Permissions{
		ac: service.DefaultAccessControl,
	}
}

func (ctrl Permissions) Effective(ctx context.Context, r *request.PermissionsEffective) (interface{}, error) {
	// Without a resource the answer is the component's own operations — the
	// ones a caller asks about before it has anything to point at.
	if r.Resource == "" {
		return ctrl.ac.Effective(ctx, types.Component{}), nil
	}

	return ctrl.ac.EffectiveFor(ctx, r.Resource)
}

func (ctrl Permissions) Trace(ctx context.Context, r *request.PermissionsTrace) (interface{}, error) {
	return ctrl.ac.Trace(ctx, r.UserID, r.RoleID, r.Resource...)
}

func (ctrl Permissions) List(ctx context.Context, r *request.PermissionsList) (interface{}, error) {
	return ctrl.ac.List(), nil
}

func (ctrl Permissions) Read(ctx context.Context, r *request.PermissionsRead) (interface{}, error) {
	return ctrl.ac.FindRules(ctx, r.RoleID, r.Resource...)
}

func (ctrl Permissions) Delete(ctx context.Context, r *request.PermissionsDelete) (interface{}, error) {
	rr, err := ctrl.ac.FindRulesByRoleID(ctx, r.RoleID)
	if err != nil {
		return nil, err
	}

	// A role's rules span every component; this endpoint, and the grant behind
	// it, reach only as far as this one. Another component's rules are for its
	// own endpoint to clear.
	component := rbac.ResourceComponent(types.ComponentRbacResource())

	own := make(rbac.RuleSet, 0, len(rr))
	for _, rule := range rr {
		if rbac.ResourceComponent(rule.Resource) != component {
			continue
		}

		rule.Access = rbac.Inherit
		own = append(own, rule)
	}

	return api.OK(), ctrl.ac.Grant(ctx, own...)
}

func (ctrl Permissions) Update(ctx context.Context, r *request.PermissionsUpdate) (interface{}, error) {
	for _, rule := range r.Rules {
		// Make sure everything is properly set
		rule.RoleID = r.RoleID
	}

	return api.OK(), ctrl.ac.Grant(ctx, r.Rules...)
}
