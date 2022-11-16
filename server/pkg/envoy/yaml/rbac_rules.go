package yaml

import (
	"strconv"
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza/server/pkg/rbac"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	rbacRule struct {
		res *rbac.Rule

		// The resource in question
		refRbacResource string
		refRbacRes      *resource.Ref

		// The path items
		refPathRes []*resource.Ref

		// The role
		refRole *resource.Ref
		role    *types.Role
	}
	rbacRuleSet []*rbacRule
)

func rbacRuleFromResource(r *resource.RbacRule, cfg *EncoderConfig) *rbacRule {
	rr := r.Res.Resource
	rr = strings.Trim(rr, ":")

	if !strings.Contains(rr, ":") {
		rr = "corteza::" + rr
	}

	return &rbacRule{
		res: r.Res,

		refRbacResource: rr,
		refRbacRes:      r.RefRes,

		refPathRes: r.RefPath,

		refRole: r.RefRole,
	}
}

func (rr rbacRuleSet) groupByAccess() []rbacRuleSet {
	rtr := make([]rbacRuleSet, 2)

	for _, r := range rr {
		if r.res.Access == rbac.Allow {
			rtr[0] = append(rtr[0], r)
		} else if r.res.Access == rbac.Deny {
			rtr[1] = append(rtr[1], r)
		}
	}

	return rtr
}

func (rr rbacRuleSet) groupByRole() []rbacRuleSet {
	rolx := make(map[string]rbacRuleSet)

	for _, r := range rr {
		roleKey := r.getRoleKey()

		if _, has := rolx[roleKey]; !has {
			rolx[roleKey] = make(rbacRuleSet, 0, 100)
		}

		rolx[roleKey] = append(rolx[roleKey], r)
	}

	rtr := make([]rbacRuleSet, 0, len(rolx))
	for _, rr := range rolx {
		rtr = append(rtr, rr)
	}

	return rtr
}

func (rr rbacRuleSet) groupByResource() []rbacRuleSet {
	rolx := make(map[string]rbacRuleSet)

	for _, r := range rr {
		k := r.res.Resource
		if _, has := rolx[k]; !has {
			rolx[k] = make(rbacRuleSet, 0, 100)
		}

		rolx[k] = append(rolx[k], r)
	}

	rtr := make([]rbacRuleSet, 0, len(rolx))
	for _, rr := range rolx {
		rtr = append(rtr, rr)
	}

	return rtr
}

func (r *rbacRule) getRoleKey() string {
	if r.role == nil {
		return ""
	}
	if r.role.Handle != "" {
		return r.role.Handle
	}
	return strconv.FormatUint(r.res.RoleID, 10)
}
